package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var projectViewCmd = &cobra.Command{
	Use:   "view <PROJECT-ID>",
	Short: "Просмотр проекта",
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
		project, err := c.GetProject(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(project)
		}

		lead := ""
		if project.Lead != nil {
			lead = format.StripControlChars(project.Lead.Name)
		}

		headers := []string{"FIELD", "VALUE"}
		rows := [][]string{
			{"ID", project.ID},
			{"Name", format.StripControlChars(project.Name)},
			{"State", project.State},
			{"Description", format.StripControlChars(project.Description)},
			{"Start Date", project.StartDate},
			{"Target Date", project.TargetDate},
			{"Lead", lead},
			{"URL", project.URL},
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать проект",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		teamIDsStr, _ := cmd.Flags().GetString("team-ids")
		startDate, _ := cmd.Flags().GetString("start-date")
		targetDate, _ := cmd.Flags().GetString("target-date")
		leadID, _ := cmd.Flags().GetString("lead-id")
		color, _ := cmd.Flags().GetString("color")
		icon, _ := cmd.Flags().GetString("icon")
		memberIDsStr, _ := cmd.Flags().GetString("member-ids")
		content, _ := cmd.Flags().GetString("content")

		if name == "" {
			return fmt.Errorf("требуется --name")
		}
		if teamIDsStr == "" {
			return fmt.Errorf("требуется --team-ids")
		}

		teamIDs := splitComma(teamIDsStr)

		var memberIDs []string
		if memberIDsStr != "" {
			memberIDs = splitComma(memberIDsStr)
		}

		inp := client.CreateProjectInput{
			Name:        name,
			Description: description,
			TeamIDs:     teamIDs,
			StartDate:   startDate,
			TargetDate:  targetDate,
			LeadID:      leadID,
			Color:       color,
			Icon:        icon,
			MemberIDs:   memberIDs,
			Content:     content,
		}
		if cmd.Flags().Changed("priority") {
			p, _ := cmd.Flags().GetInt("priority")
			inp.Priority = &p
		}

		c := newLinearClient(t)
		project, err := c.CreateProject(inp)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Проект создан: %s (%s)\n", project.Name, project.ID)
		return nil
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update <PROJECT-ID>",
	Short: "Обновить проект",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		startDate, _ := cmd.Flags().GetString("start-date")
		targetDate, _ := cmd.Flags().GetString("target-date")
		leadID, _ := cmd.Flags().GetString("lead-id")
		state, _ := cmd.Flags().GetString("state")
		color, _ := cmd.Flags().GetString("color")
		icon, _ := cmd.Flags().GetString("icon")
		memberIDsStr, _ := cmd.Flags().GetString("member-ids")
		content, _ := cmd.Flags().GetString("content")

		hasFlag := name != "" || description != "" || startDate != "" || targetDate != "" ||
			leadID != "" || state != "" || color != "" || icon != "" || memberIDsStr != "" ||
			content != "" || cmd.Flags().Changed("priority")
		if !hasFlag {
			return fmt.Errorf("укажите хотя бы один флаг для обновления: --name, --description, --start-date, --target-date, --lead-id, --state, --color, --icon, --priority, --member-ids, --content")
		}

		var memberIDs []string
		if memberIDsStr != "" {
			memberIDs = splitComma(memberIDsStr)
		}

		upd := client.UpdateProjectInput{
			Name:        name,
			Description: description,
			StartDate:   startDate,
			TargetDate:  targetDate,
			LeadID:      leadID,
			State:       state,
			Color:       color,
			Icon:        icon,
			MemberIDs:   memberIDs,
			Content:     content,
		}
		if cmd.Flags().Changed("priority") {
			p, _ := cmd.Flags().GetInt("priority")
			upd.Priority = &p
		}

		c := newLinearClient(t)
		project, err := c.UpdateProject(args[0], upd)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Проект обновлён: %s (%s)\n", project.Name, project.ID)
		return nil
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete <PROJECT-ID>",
	Short: "Удалить проект",
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
		if err := c.DeleteProject(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Проект удалён: %s\n", args[0])
		return nil
	},
}

var projectArchiveCmd = &cobra.Command{
	Use:   "archive <PROJECT-ID>",
	Short: "Архивировать проект",
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
		if err := c.ArchiveProject(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Проект архивирован: %s\n", args[0])
		return nil
	},
}

var projectUnarchiveCmd = &cobra.Command{
	Use:   "unarchive <PROJECT-ID>",
	Short: "Разархивировать проект",
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
		if err := c.UnarchiveProject(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Проект разархивирован: %s\n", args[0])
		return nil
	},
}

func init() {
	projectCreateCmd.Flags().String("name", "", "Название проекта (обязательно)")
	projectCreateCmd.Flags().String("description", "", "Описание проекта")
	projectCreateCmd.Flags().String("team-ids", "", "ID команд через запятую (обязательно)")
	projectCreateCmd.Flags().String("start-date", "", "Дата начала (YYYY-MM-DD)")
	projectCreateCmd.Flags().String("target-date", "", "Целевая дата (YYYY-MM-DD)")
	projectCreateCmd.Flags().String("lead-id", "", "ID руководителя проекта")
	projectCreateCmd.Flags().String("color", "", "Цвет проекта (hex, например #FF0000)")
	projectCreateCmd.Flags().String("icon", "", "Иконка проекта (эмодзи или имя)")
	projectCreateCmd.Flags().Int("priority", 0, "Приоритет проекта (0=нет, 1=срочный, 2=высокий, 3=средний, 4=низкий)")
	projectCreateCmd.Flags().String("member-ids", "", "ID участников через запятую")
	projectCreateCmd.Flags().String("content", "", "Содержимое/описание проекта (markdown)")

	projectUpdateCmd.Flags().String("name", "", "Новое название проекта")
	projectUpdateCmd.Flags().String("description", "", "Новое описание")
	projectUpdateCmd.Flags().String("start-date", "", "Дата начала (YYYY-MM-DD)")
	projectUpdateCmd.Flags().String("target-date", "", "Целевая дата (YYYY-MM-DD)")
	projectUpdateCmd.Flags().String("lead-id", "", "ID руководителя проекта")
	projectUpdateCmd.Flags().String("state", "", "Статус проекта")
	projectUpdateCmd.Flags().String("color", "", "Цвет проекта (hex, например #FF0000)")
	projectUpdateCmd.Flags().String("icon", "", "Иконка проекта (эмодзи или имя)")
	projectUpdateCmd.Flags().Int("priority", 0, "Приоритет проекта (0=нет, 1=срочный, 2=высокий, 3=средний, 4=низкий)")
	projectUpdateCmd.Flags().String("member-ids", "", "ID участников через запятую")
	projectUpdateCmd.Flags().String("content", "", "Содержимое/описание проекта (markdown)")

	projectsCmd.AddCommand(projectViewCmd)
	projectsCmd.AddCommand(projectCreateCmd)
	projectsCmd.AddCommand(projectUpdateCmd)
	projectsCmd.AddCommand(projectDeleteCmd)
	projectsCmd.AddCommand(projectArchiveCmd)
	projectsCmd.AddCommand(projectUnarchiveCmd)
}
