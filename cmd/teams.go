package cmd

import (
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
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		teams, err := c.ListTeams()
		if err != nil {
			return err
		}

		if len(teams) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Команды не найдены.")
			return nil
		}

		headers := []string{"ID", "KEY", "NAME"}
		rows := make([][]string, 0, len(teams))
		for _, team := range teams {
			rows = append(rows, []string{team.ID, team.Key, team.Name})
		}
		format.TableWriter(cmd.OutOrStdout(), headers, rows)
		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamsListCmd)
	rootCmd.AddCommand(teamsCmd)
}
