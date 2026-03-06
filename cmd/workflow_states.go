package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var workflowStatesCmd = &cobra.Command{
	Use:   "workflow-states",
	Short: "Управление статусами задач (workflow states)",
}

var workflowStatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список статусов задач команды",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		teamKey, _ := cmd.Flags().GetString("team")
		if teamKey == "" {
			return fmt.Errorf("требуется --team (ключ или ID команды)")
		}

		c := newLinearClient(t)
		teamID := teamKey
		if team, err := c.GetTeamByKey(teamKey); err != nil {
			return fmt.Errorf("ошибка получения команды: %w", err)
		} else if team != nil {
			teamID = team.ID
		}
		states, err := c.ListWorkflowStates(teamID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(states)
		}

		if len(states) == 0 {
			fmt.Fprintln(out, "Статусы не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "TYPE", "COLOR", "TEAM"}
		rows := make([][]string, 0, len(states))
		for _, s := range states {
			rows = append(rows, []string{
				s.ID,
				format.StripControlChars(s.Name),
				s.Type,
				s.Color,
				format.StripControlChars(s.Team.Name),
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var workflowStatesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новый статус задачи",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		teamKey, _ := cmd.Flags().GetString("team")
		if teamKey == "" {
			return fmt.Errorf("требуется --team (ключ или ID команды)")
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("требуется --name (название статуса)")
		}
		stateType, _ := cmd.Flags().GetString("type")
		if stateType == "" {
			return fmt.Errorf("требуется --type (тип: triage, backlog, unstarted, started, completed, cancelled)")
		}
		color, _ := cmd.Flags().GetString("color")

		c := newLinearClient(t)
		teamID := teamKey
		if team, err := c.GetTeamByKey(teamKey); err != nil {
			return fmt.Errorf("ошибка получения команды: %w", err)
		} else if team != nil {
			teamID = team.ID
		}
		ws, err := c.CreateWorkflowState(teamID, name, stateType, color)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(ws)
		}

		fmt.Fprintf(out, "Статус создан: %s (ID: %s)\n", format.StripControlChars(ws.Name), ws.ID)
		return nil
	},
}

var workflowStatesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить статус задачи",
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
		color, _ := cmd.Flags().GetString("color")

		c := newLinearClient(t)
		ws, err := c.UpdateWorkflowState(args[0], name, color)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(ws)
		}

		fmt.Fprintf(out, "Статус обновлён: %s (ID: %s)\n", format.StripControlChars(ws.Name), ws.ID)
		return nil
	},
}

var workflowStatesArchiveCmd = &cobra.Command{
	Use:   "archive <ID>",
	Short: "Архивировать статус задачи",
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
		if err := c.ArchiveWorkflowState(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Статус %s архивирован.\n", args[0])
		return nil
	},
}

func init() {
	workflowStatesListCmd.Flags().String("team", "", "ID команды (обязательно)")

	workflowStatesCreateCmd.Flags().String("team", "", "ID команды (обязательно)")
	workflowStatesCreateCmd.Flags().String("name", "", "Название статуса (обязательно)")
	workflowStatesCreateCmd.Flags().String("type", "", "Тип статуса: triage, backlog, unstarted, started, completed, cancelled (обязательно)")
	workflowStatesCreateCmd.Flags().String("color", "", "Цвет статуса (hex, например: #e2e2e2)")

	workflowStatesUpdateCmd.Flags().String("name", "", "Новое название статуса")
	workflowStatesUpdateCmd.Flags().String("color", "", "Новый цвет статуса")

	workflowStatesCmd.AddCommand(workflowStatesListCmd)
	workflowStatesCmd.AddCommand(workflowStatesCreateCmd)
	workflowStatesCmd.AddCommand(workflowStatesUpdateCmd)
	workflowStatesCmd.AddCommand(workflowStatesArchiveCmd)
	rootCmd.AddCommand(workflowStatesCmd)
}
