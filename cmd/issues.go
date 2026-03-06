package cmd

import (
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

// newLinearClient — фабрика клиента, заменяется в тестах.
var newLinearClient = func(token string) *client.Client {
	return client.New(token)
}

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Управление задачами",
}

var (
	issuesAssignee string
	issuesStatus   string
	issuesTeam     string
	issuesLimit    int
)

var issuesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список задач",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		issues, err := c.ListIssues(client.IssueFilter{
			Assignee: issuesAssignee,
			Status:   issuesStatus,
			Team:     issuesTeam,
			Limit:    issuesLimit,
		})
		if err != nil {
			return err
		}

		if len(issues) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Задачи не найдены.")
			return nil
		}

		headers := []string{"ID", "TITLE", "STATUS", "ASSIGNEE", "PRIORITY", "UPDATED"}
		rows := make([][]string, 0, len(issues))
		for _, issue := range issues {
			assignee := ""
			if issue.Assignee != nil {
				assignee = issue.Assignee.DisplayName
			}
			rows = append(rows, []string{
				issue.Identifier,
				issue.Title,
				issue.State.Name,
				assignee,
				client.PriorityLabel(issue.Priority),
				issue.UpdatedAt.Format("2006-01-02"),
			})
		}
		format.TableWriter(cmd.OutOrStdout(), headers, rows)
		return nil
	},
}

func init() {
	issuesListCmd.Flags().StringVar(&issuesAssignee, "assignee", "", "Фильтр по исполнителю")
	issuesListCmd.Flags().StringVar(&issuesStatus, "status", "", "Фильтр по статусу")
	issuesListCmd.Flags().StringVar(&issuesTeam, "team", "", "Фильтр по команде (key)")
	issuesListCmd.Flags().IntVar(&issuesLimit, "limit", 25, "Максимальное количество задач")

	issuesCmd.AddCommand(issuesListCmd)
	rootCmd.AddCommand(issuesCmd)
}
