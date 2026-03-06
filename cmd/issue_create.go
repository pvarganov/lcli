package cmd

import (
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/spf13/cobra"
)

var (
	createTitle       string
	createDescription string
	createTeam        string
	createAssignee    string
	createPriority    int
)

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новую задачу",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)

		team, err := c.GetTeamByKey(createTeam)
		if err != nil {
			return fmt.Errorf("ошибка получения команды: %w", err)
		}
		if team == nil {
			return fmt.Errorf("команда %q не найдена", createTeam)
		}

		input := client.CreateIssueInput{
			TeamID:      team.ID,
			Title:       createTitle,
			Description: createDescription,
			Priority:    createPriority,
		}

		if createAssignee != "" {
			user, err := c.FindUserByName(createAssignee)
			if err != nil {
				return fmt.Errorf("ошибка поиска пользователя: %w", err)
			}
			if user == nil {
				return fmt.Errorf("пользователь %q не найден", createAssignee)
			}
			input.AssigneeID = user.ID
		}

		issue, err := c.CreateIssue(input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Задача создана: %s  %s\n", issue.Identifier, issue.Title)
		return nil
	},
}

func init() {
	issueCreateCmd.Flags().StringVar(&createTitle, "title", "", "Заголовок задачи (обязательно)")
	issueCreateCmd.Flags().StringVar(&createDescription, "description", "", "Описание задачи")
	issueCreateCmd.Flags().StringVar(&createTeam, "team", "", "Ключ команды (обязательно)")
	issueCreateCmd.Flags().StringVar(&createAssignee, "assignee", "", "Исполнитель (displayName или email)")
	issueCreateCmd.Flags().IntVar(&createPriority, "priority", 0, "Приоритет: 1=Urgent, 2=High, 3=Medium, 4=Low")
	_ = issueCreateCmd.MarkFlagRequired("title")
	_ = issueCreateCmd.MarkFlagRequired("team")

	issueCmd.AddCommand(issueCreateCmd)
}
