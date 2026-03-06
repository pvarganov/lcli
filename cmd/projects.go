package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Управление проектами",
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список проектов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		projects, err := c.ListProjects()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(projects)
		}

		if len(projects) == 0 {
			fmt.Fprintln(out, "Проекты не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "STATE", "DESCRIPTION"}
		rows := make([][]string, 0, len(projects))
		for _, p := range projects {
			desc := p.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			rows = append(rows, []string{p.ID, p.Name, p.State, desc})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

func init() {
	projectsCmd.AddCommand(projectsListCmd)
	rootCmd.AddCommand(projectsCmd)
}
