package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Управление шаблонами",
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список шаблонов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		templates, err := c.ListTemplates()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(templates)
		}

		if len(templates) == 0 {
			fmt.Fprintln(out, "Шаблоны не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "TYPE", "TEAM"}
		rows := make([][]string, 0, len(templates))
		for _, tpl := range templates {
			teamName := ""
			if tpl.Team != nil {
				teamName = tpl.Team.Name
			}
			rows = append(rows, []string{tpl.ID, tpl.Name, tpl.Type, teamName})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var templatesViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр шаблона по ID",
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
		tpl, err := c.GetTemplate(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(tpl)
		}

		fmt.Fprintf(out, "ID:          %s\n", tpl.ID)
		fmt.Fprintf(out, "Название:    %s\n", tpl.Name)
		fmt.Fprintf(out, "Тип:         %s\n", tpl.Type)
		if tpl.Description != "" {
			fmt.Fprintf(out, "Описание:    %s\n", tpl.Description)
		}
		if tpl.Team != nil {
			fmt.Fprintf(out, "Команда:     %s\n", tpl.Team.Name)
		}
		if tpl.Creator != nil {
			fmt.Fprintf(out, "Автор:       %s\n", tpl.Creator.Name)
		}
		return nil
	},
}

var templatesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать шаблон",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		templateType, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		teamID, _ := cmd.Flags().GetString("team")

		if name == "" {
			return fmt.Errorf("необходимо указать --name")
		}
		if templateType == "" {
			return fmt.Errorf("необходимо указать --type")
		}

		input := map[string]any{}
		if description != "" {
			input["description"] = description
		}
		if teamID != "" {
			input["teamId"] = teamID
		}

		c := newLinearClient(t)
		tpl, err := c.CreateTemplate(name, templateType, input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(tpl)
		}
		fmt.Fprintf(out, "Шаблон создан: %s (%s)\n", tpl.Name, tpl.ID)
		return nil
	},
}

var templatesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить шаблон по ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		input := map[string]any{}
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			input["name"] = name
		}
		if description, _ := cmd.Flags().GetString("description"); description != "" {
			input["description"] = description
		}

		if len(input) == 0 {
			return fmt.Errorf("необходимо указать хотя бы один флаг для обновления")
		}

		c := newLinearClient(t)
		tpl, err := c.UpdateTemplate(args[0], input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(tpl)
		}
		fmt.Fprintf(out, "Шаблон обновлён: %s (%s)\n", tpl.Name, tpl.ID)
		return nil
	},
}

var templatesDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить шаблон по ID",
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
		if err := c.DeleteTemplate(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Шаблон удалён: %s\n", args[0])
		return nil
	},
}

func init() {
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesViewCmd)

	templatesCreateCmd.Flags().String("name", "", "Название шаблона (обязательно)")
	templatesCreateCmd.Flags().String("type", "", "Тип шаблона (обязательно), например: issue")
	templatesCreateCmd.Flags().String("description", "", "Описание шаблона")
	templatesCreateCmd.Flags().String("team", "", "ID команды")
	templatesCmd.AddCommand(templatesCreateCmd)

	templatesUpdateCmd.Flags().String("name", "", "Новое название")
	templatesUpdateCmd.Flags().String("description", "", "Новое описание")
	templatesCmd.AddCommand(templatesUpdateCmd)

	templatesCmd.AddCommand(templatesDeleteCmd)

	rootCmd.AddCommand(templatesCmd)
}
