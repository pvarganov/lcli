package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var (
	batchUpdateIDs    string
	batchUpdateStatus string
)

var issueBatchUpdateCmd = &cobra.Command{
	Use:   "batch-update",
	Short: "Обновить несколько задач за один запрос",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}
		if batchUpdateIDs == "" {
			return fmt.Errorf("укажите --ids (через запятую)")
		}
		if batchUpdateStatus == "" {
			return fmt.Errorf("укажите --status")
		}

		ids := strings.Split(batchUpdateIDs, ",")
		for i, id := range ids {
			ids[i] = strings.TrimSpace(id)
		}

		c := newLinearClient(t)

		// Найти state ID по имени — берём через первую задачу в списке
		// Сначала получим команду из первой задачи для поиска состояния
		stateID, err := findStateIDForBatch(c, ids[0], batchUpdateStatus)
		if err != nil {
			return err
		}
		if stateID == "" {
			return fmt.Errorf("статус не найден: %s", batchUpdateStatus)
		}

		issues, err := c.BatchUpdateIssues(client.BatchUpdateIssuesInput{
			IDs:    ids,
			Update: client.UpdateIssueInput{StateID: stateID},
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

		fmt.Fprintf(out, "Обновлено задач: %d\n", len(issues))
		if len(issues) > 0 {
			headers := []string{"ID", "TITLE", "STATUS"}
			rows := make([][]string, 0, len(issues))
			for _, iss := range issues {
				rows = append(rows, []string{
					iss.Identifier,
					format.StripControlChars(iss.Title),
					format.ColorStatus(iss.State.Name, iss.State.Type),
				})
			}
			format.TableWriter(out, headers, rows)
		}
		return nil
	},
}

// findStateIDForBatch получает ID состояния задачи по её идентификатору и имени статуса.
func findStateIDForBatch(c *client.Client, issueID, stateName string) (string, error) {
	issue, err := c.GetIssue(issueID)
	if err != nil {
		return "", err
	}
	if issue == nil {
		return "", fmt.Errorf("задача не найдена: %s", issueID)
	}
	return c.FindWorkflowStateByName(issue.Team.ID, stateName)
}

func init() {
	issueBatchUpdateCmd.Flags().StringVar(&batchUpdateIDs, "ids", "", "ID задач через запятую (обязательно)")
	issueBatchUpdateCmd.Flags().StringVar(&batchUpdateStatus, "status", "", "Новый статус (обязательно)")

	issuesCmd.AddCommand(issueBatchUpdateCmd)
}
