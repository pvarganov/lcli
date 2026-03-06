package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Работа с отдельной задачей",
}

var issueViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Показать детали задачи",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		issue, err := c.GetIssue(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(issue)
		}

		fmt.Fprintf(out, "%s  %s\n", issue.Identifier, issue.Title)
		fmt.Fprintln(out, strings.Repeat("─", 60))
		fmt.Fprintf(out, "Статус:     %s\n", issue.State.Name)
		fmt.Fprintf(out, "Приоритет:  %s\n", client.PriorityLabel(issue.Priority))
		if issue.Assignee != nil {
			fmt.Fprintf(out, "Исполнитель: %s\n", issue.Assignee.DisplayName)
		} else {
			fmt.Fprintln(out, "Исполнитель: —")
		}
		fmt.Fprintf(out, "Команда:    %s\n", issue.Team.Name)
		fmt.Fprintf(out, "Обновлено:  %s\n", issue.UpdatedAt.Format("2006-01-02 15:04"))
		if issue.Description != "" {
			fmt.Fprintln(out, strings.Repeat("─", 60))
			fmt.Fprintln(out, issue.Description)
		}
		return nil
	},
}

func init() {
	issueCmd.AddCommand(issueViewCmd)
	rootCmd.AddCommand(issueCmd)
}
