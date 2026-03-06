package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var projectStatusesCmd = &cobra.Command{
	Use:   "statuses",
	Short: "Управление статусами проектов",
}

var projectStatusesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список статусов проектов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		statuses, err := c.ListProjectStatuses()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(statuses)
		}

		if len(statuses) == 0 {
			fmt.Fprintln(out, "Статусы проектов не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "TYPE", "COLOR", "DESCRIPTION"}
		rows := make([][]string, 0, len(statuses))
		for _, s := range statuses {
			desc := format.StripControlChars(s.Description)
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			rows = append(rows, []string{s.ID, format.StripControlChars(s.Name), s.Type, s.Color, desc})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var projectStatusesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать статус проекта",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		statusType, _ := cmd.Flags().GetString("type")
		color, _ := cmd.Flags().GetString("color")
		description, _ := cmd.Flags().GetString("description")
		position, _ := cmd.Flags().GetFloat64("position")

		if name == "" {
			return fmt.Errorf("требуется --name")
		}
		if statusType == "" {
			return fmt.Errorf("требуется --type (backlog, planned, started, paused, completed, canceled)")
		}
		if color == "" {
			return fmt.Errorf("требуется --color")
		}

		c := newLinearClient(t)
		status, err := c.CreateProjectStatus(client.CreateProjectStatusInput{
			Name:        name,
			Type:        statusType,
			Color:       color,
			Description: description,
			Position:    position,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Статус создан: %s (%s)\n", status.Name, status.ID)
		return nil
	},
}

var projectStatusesUpdateCmd = &cobra.Command{
	Use:   "update <STATUS-ID>",
	Short: "Обновить статус проекта",
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
		description, _ := cmd.Flags().GetString("description")

		c := newLinearClient(t)
		status, err := c.UpdateProjectStatus(args[0], client.UpdateProjectStatusInput{
			Name:        name,
			Color:       color,
			Description: description,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Статус обновлён: %s (%s)\n", status.Name, status.ID)
		return nil
	},
}

func init() {
	projectStatusesCreateCmd.Flags().String("name", "", "Название статуса (обязательно)")
	projectStatusesCreateCmd.Flags().String("type", "", "Тип статуса: backlog, planned, started, paused, completed, canceled (обязательно)")
	projectStatusesCreateCmd.Flags().String("color", "", "Цвет статуса (hex, например #ff0000) (обязательно)")
	projectStatusesCreateCmd.Flags().String("description", "", "Описание статуса")
	projectStatusesCreateCmd.Flags().Float64("position", 0, "Позиция статуса (порядок)")

	projectStatusesUpdateCmd.Flags().String("name", "", "Новое название статуса")
	projectStatusesUpdateCmd.Flags().String("color", "", "Новый цвет статуса")
	projectStatusesUpdateCmd.Flags().String("description", "", "Новое описание статуса")

	projectStatusesCmd.AddCommand(projectStatusesListCmd)
	projectStatusesCmd.AddCommand(projectStatusesCreateCmd)
	projectStatusesCmd.AddCommand(projectStatusesUpdateCmd)
	projectsCmd.AddCommand(projectStatusesCmd)
}
