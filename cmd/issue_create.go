package cmd

import (
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var (
	createTitle       string
	createDescription string
	createTeam        string
	createAssignee    string
	createPriority    int
	createDueDate     string
	createEstimate    int
	createLabels      string
	createParent      string
	createState       string
	createCycleID     string
	createProjectID   string
	createMilestoneID string
)

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новую задачу",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
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
			DueDate:     createDueDate,
			CycleID:     createCycleID,
			ProjectID:   createProjectID,
			MilestoneID: createMilestoneID,
			ParentID:    createParent,
		}

		if cmd.Flags().Changed("estimate") {
			input.Estimate = &createEstimate
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

		if createLabels != "" {
			names := splitComma(createLabels)
			ids, err := c.FindLabelsByNames(team.ID, names)
			if err != nil {
				return fmt.Errorf("ошибка поиска меток: %w", err)
			}
			input.LabelIDs = ids
		}

		if createState != "" {
			stateID, err := c.FindWorkflowStateByName(team.ID, createState)
			if err != nil {
				return fmt.Errorf("ошибка поиска состояния: %w", err)
			}
			if stateID == "" {
				return fmt.Errorf("состояние %q не найдено", createState)
			}
			input.StateID = stateID
		}

		issue, err := c.CreateIssue(input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Задача создана: %s  %s\n", issue.Identifier, format.StripControlChars(issue.Title))
		return nil
	},
}

func init() {
	issueCreateCmd.Flags().StringVar(&createTitle, "title", "", "Заголовок задачи (обязательно)")
	issueCreateCmd.Flags().StringVar(&createDescription, "description", "", "Описание задачи")
	issueCreateCmd.Flags().StringVar(&createTeam, "team", "", "Ключ команды (обязательно)")
	issueCreateCmd.Flags().StringVar(&createAssignee, "assignee", "", "Исполнитель (displayName или email)")
	issueCreateCmd.Flags().IntVar(&createPriority, "priority", 0, "Приоритет: 1=Urgent, 2=High, 3=Medium, 4=Low")
	issueCreateCmd.Flags().StringVar(&createDueDate, "due-date", "", "Срок выполнения (YYYY-MM-DD)")
	issueCreateCmd.Flags().IntVar(&createEstimate, "estimate", 0, "Оценка задачи в story points")
	issueCreateCmd.Flags().StringVar(&createLabels, "labels", "", "Метки через запятую (по имени)")
	issueCreateCmd.Flags().StringVar(&createParent, "parent", "", "ID родительской задачи")
	issueCreateCmd.Flags().StringVar(&createState, "state", "", "Состояние задачи (по имени)")
	issueCreateCmd.Flags().StringVar(&createCycleID, "cycle-id", "", "ID цикла")
	issueCreateCmd.Flags().StringVar(&createProjectID, "project-id", "", "ID проекта")
	issueCreateCmd.Flags().StringVar(&createMilestoneID, "milestone-id", "", "ID вехи проекта")
	_ = issueCreateCmd.MarkFlagRequired("title")
	_ = issueCreateCmd.MarkFlagRequired("team")

	issueCmd.AddCommand(issueCreateCmd)
}
