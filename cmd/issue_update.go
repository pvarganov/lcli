package cmd

import (
	"fmt"
	"strings"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var (
	updateStatus      string
	updateAssignee    string
	updatePriority    int
	updateTitle       string
	updateDescription string
	updateDueDate     string
	updateEstimate    int
	updateParent      string
	updateCycleID     string
	updateProjectID   string
	updateMilestoneID string
	updateAddLabels   string
	updateRemoveLabels string
	updateSnoozeUntil string
)

var issueUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		// Проверяем, что задан хотя бы один флаг для обновления
		hasFlag := updateTitle != "" || updateStatus != "" || updateAssignee != "" ||
			cmd.Flags().Changed("priority") || updateDescription != "" || updateDueDate != "" ||
			cmd.Flags().Changed("estimate") || updateParent != "" || updateCycleID != "" ||
			updateProjectID != "" || updateMilestoneID != "" || updateAddLabels != "" ||
			updateRemoveLabels != "" || updateSnoozeUntil != ""
		if !hasFlag {
			return fmt.Errorf("укажите хотя бы один флаг для обновления: --title, --status, --assignee, --priority, --description, --due-date, --estimate, --parent, --cycle-id, --project-id, --milestone-id, --add-labels, --remove-labels или --snooze-until")
		}

		issueID := args[0]
		c := newLinearClient(t)

		input := client.UpdateIssueInput{
			Title:          updateTitle,
			Description:    updateDescription,
			DueDate:        updateDueDate,
			ParentID:       updateParent,
			CycleID:        updateCycleID,
			ProjectID:      updateProjectID,
			MilestoneID:    updateMilestoneID,
			SnoozedUntilAt: updateSnoozeUntil,
		}

		if updateAssignee != "" {
			user, err := c.FindUserByName(updateAssignee)
			if err != nil {
				return fmt.Errorf("ошибка поиска пользователя: %w", err)
			}
			if user == nil {
				return fmt.Errorf("пользователь %q не найден", updateAssignee)
			}
			input.AssigneeID = user.ID
		}

		if updateStatus != "" {
			issue, err := c.GetIssue(issueID)
			if err != nil {
				return fmt.Errorf("ошибка получения задачи: %w", err)
			}
			stateID, err := c.FindWorkflowStateByName(issue.Team.ID, updateStatus)
			if err != nil {
				return fmt.Errorf("ошибка поиска статуса: %w", err)
			}
			if stateID == "" {
				return fmt.Errorf("статус %q не найден в команде %s", updateStatus, issue.Team.Name)
			}
			input.StateID = stateID
		}

		if cmd.Flags().Changed("priority") {
			p := updatePriority
			input.Priority = &p
		}

		if cmd.Flags().Changed("estimate") {
			e := updateEstimate
			input.Estimate = &e
		}

		if updateAddLabels != "" || updateRemoveLabels != "" {
			issue, err := c.GetIssue(issueID)
			if err != nil {
				return fmt.Errorf("ошибка получения задачи: %w", err)
			}
			teamID := issue.Team.ID

			if updateAddLabels != "" {
				names := splitComma(updateAddLabels)
				ids, err := c.FindLabelsByNames(teamID, names)
				if err != nil {
					return fmt.Errorf("ошибка поиска меток для добавления: %w", err)
				}
				input.AddedLabelIDs = ids
			}

			if updateRemoveLabels != "" {
				names := splitComma(updateRemoveLabels)
				ids, err := c.FindLabelsByNames(teamID, names)
				if err != nil {
					return fmt.Errorf("ошибка поиска меток для удаления: %w", err)
				}
				input.RemovedLabelIDs = ids
			}
		}

		issue, err := c.UpdateIssue(issueID, input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Задача обновлена: %s  %s\n", issue.Identifier, format.StripControlChars(issue.Title))
		return nil
	},
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func init() {
	issueUpdateCmd.Flags().StringVar(&updateStatus, "status", "", "Новый статус задачи")
	issueUpdateCmd.Flags().StringVar(&updateAssignee, "assignee", "", "Исполнитель (displayName или email)")
	issueUpdateCmd.Flags().IntVar(&updatePriority, "priority", 0, "Приоритет: 1=Urgent, 2=High, 3=Medium, 4=Low")
	issueUpdateCmd.Flags().StringVar(&updateTitle, "title", "", "Новый заголовок задачи")
	issueUpdateCmd.Flags().StringVar(&updateDescription, "description", "", "Описание задачи (Markdown)")
	issueUpdateCmd.Flags().StringVar(&updateDueDate, "due-date", "", "Срок выполнения (YYYY-MM-DD)")
	issueUpdateCmd.Flags().IntVar(&updateEstimate, "estimate", 0, "Оценка задачи в story points")
	issueUpdateCmd.Flags().StringVar(&updateParent, "parent", "", "ID родительской задачи")
	issueUpdateCmd.Flags().StringVar(&updateCycleID, "cycle-id", "", "ID цикла (спринта)")
	issueUpdateCmd.Flags().StringVar(&updateProjectID, "project-id", "", "ID проекта")
	issueUpdateCmd.Flags().StringVar(&updateMilestoneID, "milestone-id", "", "ID milestone проекта")
	issueUpdateCmd.Flags().StringVar(&updateAddLabels, "add-labels", "", "Метки для добавления (через запятую)")
	issueUpdateCmd.Flags().StringVar(&updateRemoveLabels, "remove-labels", "", "Метки для удаления (через запятую)")
	issueUpdateCmd.Flags().StringVar(&updateSnoozeUntil, "snooze-until", "", "Отложить до (RFC3339, напр. 2026-04-01T10:00:00Z)")

	issueCmd.AddCommand(issueUpdateCmd)
}
