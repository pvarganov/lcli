package cmd

import (
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/spf13/cobra"
)

var (
	updateStatus   string
	updateAssignee string
	updatePriority int
	updateTitle    string
)

var issueUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		// Проверяем, что задан хотя бы один флаг для обновления
		if updateTitle == "" && updateStatus == "" && updateAssignee == "" && !cmd.Flags().Changed("priority") {
			return fmt.Errorf("укажите хотя бы один флаг для обновления: --title, --status, --assignee или --priority")
		}

		issueID := args[0]
		c := newLinearClient(t)

		input := client.UpdateIssueInput{
			Title: updateTitle,
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

		issue, err := c.UpdateIssue(issueID, input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Задача обновлена: %s  %s\n", issue.Identifier, issue.Title)
		return nil
	},
}

func init() {
	issueUpdateCmd.Flags().StringVar(&updateStatus, "status", "", "Новый статус задачи")
	issueUpdateCmd.Flags().StringVar(&updateAssignee, "assignee", "", "Исполнитель (displayName или email)")
	issueUpdateCmd.Flags().IntVar(&updatePriority, "priority", 0, "Приоритет: 1=Urgent, 2=High, 3=Medium, 4=Low")
	issueUpdateCmd.Flags().StringVar(&updateTitle, "title", "", "Новый заголовок задачи")

	issueCmd.AddCommand(issueUpdateCmd)
}
