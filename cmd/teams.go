package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Управление командами",
}

var teamsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список команд",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		teams, err := c.ListTeams()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(teams)
		}

		if len(teams) == 0 {
			fmt.Fprintln(out, "Команды не найдены.")
			return nil
		}

		headers := []string{"ID", "KEY", "NAME"}
		rows := make([][]string, 0, len(teams))
		for _, team := range teams {
			rows = append(rows, []string{team.ID, format.StripControlChars(team.Key), format.StripControlChars(team.Name)})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamsListCmd)
	rootCmd.AddCommand(teamsCmd)
}
