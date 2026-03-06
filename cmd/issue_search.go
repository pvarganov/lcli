package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var (
	issueSearchLimit int
	issueSearchTeam  string
)

var issueSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Поиск задач по строке запроса",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		issues, pageInfo, err := c.SearchIssues(args[0], issueSearchLimit)
		if err != nil {
			return err
		}

		// фильтрация по команде (client-side, если указан --team)
		if issueSearchTeam != "" {
			filtered := issues[:0]
			for _, iss := range issues {
				if iss.Team.Key == issueSearchTeam {
					filtered = append(filtered, iss)
				}
			}
			issues = filtered
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
			fmt.Fprintf(out, "\nЕсть ещё результаты, уточните запрос или увеличьте --limit\n")
		}
		return nil
	},
}

func init() {
	issueSearchCmd.Flags().IntVar(&issueSearchLimit, "limit", 25, "Максимальное количество результатов")
	issueSearchCmd.Flags().StringVar(&issueSearchTeam, "team", "", "Фильтр по команде (key)")

	issuesCmd.AddCommand(issueSearchCmd)
}
