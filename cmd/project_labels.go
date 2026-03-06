package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var projectLabelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Управление метками проектов",
}

var projectLabelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список меток проектов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		labels, err := c.ListProjectLabels()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(labels)
		}

		if len(labels) == 0 {
			fmt.Fprintln(out, "Метки проектов не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "COLOR", "DESCRIPTION"}
		rows := make([][]string, 0, len(labels))
		for _, l := range labels {
			desc := format.StripControlChars(l.Description)
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			rows = append(rows, []string{l.ID, format.StripControlChars(l.Name), l.Color, desc})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var projectLabelsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать метку проекта",
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

		if name == "" {
			return fmt.Errorf("требуется --name")
		}

		c := newLinearClient(t)
		label, err := c.CreateProjectLabel(client.CreateProjectLabelInput{
			Name:        name,
			Color:       color,
			Description: description,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Метка создана: %s (%s)\n", label.Name, label.ID)
		return nil
	},
}

var projectLabelsUpdateCmd = &cobra.Command{
	Use:   "update <LABEL-ID>",
	Short: "Обновить метку проекта",
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
		label, err := c.UpdateProjectLabel(args[0], client.UpdateProjectLabelInput{
			Name:        name,
			Color:       color,
			Description: description,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Метка обновлена: %s (%s)\n", label.Name, label.ID)
		return nil
	},
}

var projectLabelsDeleteCmd = &cobra.Command{
	Use:   "delete <LABEL-ID>",
	Short: "Удалить метку проекта",
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
		if err := c.DeleteProjectLabel(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Метка удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	projectLabelsCreateCmd.Flags().String("name", "", "Название метки (обязательно)")
	projectLabelsCreateCmd.Flags().String("color", "", "Цвет метки (hex, например #ff0000)")
	projectLabelsCreateCmd.Flags().String("description", "", "Описание метки")

	projectLabelsUpdateCmd.Flags().String("name", "", "Новое название метки")
	projectLabelsUpdateCmd.Flags().String("color", "", "Новый цвет метки")
	projectLabelsUpdateCmd.Flags().String("description", "", "Новое описание метки")

	projectLabelsCmd.AddCommand(projectLabelsListCmd)
	projectLabelsCmd.AddCommand(projectLabelsCreateCmd)
	projectLabelsCmd.AddCommand(projectLabelsUpdateCmd)
	projectLabelsCmd.AddCommand(projectLabelsDeleteCmd)
	projectsCmd.AddCommand(projectLabelsCmd)
}
