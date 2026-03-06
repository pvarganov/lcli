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

		identifiers := strings.Split(batchUpdateIDs, ",")
		for i, id := range identifiers {
			identifiers[i] = strings.TrimSpace(id)
		}
		if len(identifiers) > 50 {
			return fmt.Errorf("слишком много задач: %d (максимум 50 за один запрос)", len(identifiers))
		}

		c := newLinearClient(t)

		// Разрешаем идентификаторы (ENG-1) в UUID, требуемые issueBatchUpdate(ids: [UUID!]!)
		uuids := make([]string, 0, len(identifiers))
		var firstTeamID string
		for _, identifier := range identifiers {
			issue, err := c.GetIssue(identifier)
			if err != nil {
				return fmt.Errorf("задача не найдена: %s: %w", identifier, err)
			}
			uuids = append(uuids, issue.ID)
			if firstTeamID == "" {
				firstTeamID = issue.Team.ID
			} else if issue.Team.ID != firstTeamID {
				return fmt.Errorf("все задачи в batch-update должны принадлежать одной команде (смешение команд не поддерживается при указании --status)")
			}
		}

		stateID, err := c.FindWorkflowStateByName(firstTeamID, batchUpdateStatus)
		if err != nil {
			return err
		}
		if stateID == "" {
			return fmt.Errorf("статус не найден: %s", batchUpdateStatus)
		}

		issues, err := c.BatchUpdateIssues(client.BatchUpdateIssuesInput{
			IDs:    uuids,
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

func init() {
	issueBatchUpdateCmd.Flags().StringVar(&batchUpdateIDs, "ids", "", "ID задач через запятую (обязательно)")
	issueBatchUpdateCmd.Flags().StringVar(&batchUpdateStatus, "status", "", "Новый статус (обязательно)")

	issuesCmd.AddCommand(issueBatchUpdateCmd)
}
