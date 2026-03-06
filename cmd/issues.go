package cmd

import (
	"encoding/json"
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
	issuesAssignee  string
	issuesStatus    string
	issuesTeam      string
	issuesLimit     int
	issuesAfter     string
	issuesPriority  int
	issuesLabel     string
	issuesProjectID string
	issuesCycleID   string
	issuesCreator   string
	issuesOrderBy   string
)

var issuesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список задач",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		issues, pageInfo, err := c.ListIssues(client.IssueFilter{
			Assignee:  issuesAssignee,
			Status:    issuesStatus,
			Team:      issuesTeam,
			Limit:     issuesLimit,
			After:     issuesAfter,
			Priority:  issuesPriority,
			Label:     issuesLabel,
			ProjectID: issuesProjectID,
			CycleID:   issuesCycleID,
			Creator:   issuesCreator,
			OrderBy:   issuesOrderBy,
		})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(issues)
		}

		if len(issues) == 0 {
			fmt.Fprintln(out, "Задачи не найдены.")
			return nil
		}

		headers := []string{"ID", "TITLE", "STATUS", "ASSIGNEE", "PRIORITY", "UPDATED"}
		rows := make([][]string, 0, len(issues))
		for _, issue := range issues {
			assignee := ""
			if issue.Assignee != nil {
				assignee = format.StripControlChars(issue.Assignee.DisplayName)
			}
			rows = append(rows, []string{
				issue.Identifier,
				format.StripControlChars(issue.Title),
				format.ColorStatus(issue.State.Name, issue.State.Type),
				assignee,
				client.PriorityLabel(issue.Priority),
				issue.UpdatedAt.Format("2006-01-02"),
			})
		}
		format.TableWriter(out, headers, rows)

		if pageInfo != nil && pageInfo.HasNextPage {
			fmt.Fprintf(out, "\nСледующая страница: --after %s\n", pageInfo.EndCursor)
		}
		return nil
	},
}

func init() {
	issuesListCmd.Flags().StringVar(&issuesAssignee, "assignee", "", "Фильтр по исполнителю")
	issuesListCmd.Flags().StringVar(&issuesStatus, "status", "", "Фильтр по статусу")
	issuesListCmd.Flags().StringVar(&issuesTeam, "team", "", "Фильтр по команде (key)")
	issuesListCmd.Flags().IntVar(&issuesLimit, "limit", 25, "Максимальное количество задач")
	issuesListCmd.Flags().StringVar(&issuesAfter, "after", "", "Курсор для пагинации (из предыдущего запроса)")
	issuesListCmd.Flags().IntVar(&issuesPriority, "priority", -1, "Фильтр по приоритету (0=нет, 1=срочно, 2=высокий, 3=средний, 4=низкий)")
	issuesListCmd.Flags().StringVar(&issuesLabel, "label", "", "Фильтр по метке (имя)")
	issuesListCmd.Flags().StringVar(&issuesProjectID, "project-id", "", "Фильтр по ID проекта")
	issuesListCmd.Flags().StringVar(&issuesCycleID, "cycle-id", "", "Фильтр по ID цикла")
	issuesListCmd.Flags().StringVar(&issuesCreator, "creator", "", "Фильтр по создателю (displayName)")
	issuesListCmd.Flags().StringVar(&issuesOrderBy, "order-by", "", "Сортировка: updatedAt, createdAt, manualOrder")

	issuesCmd.AddCommand(issuesListCmd)
	rootCmd.AddCommand(issuesCmd)
}
