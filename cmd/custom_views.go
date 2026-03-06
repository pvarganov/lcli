package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var viewsCmd = &cobra.Command{
	Use:   "views",
	Short: "Управление пользовательскими представлениями",
}

var viewsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список пользовательских представлений",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		views, err := c.ListCustomViews()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(views)
		}

		if len(views) == 0 {
			fmt.Fprintln(out, "Представления не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "DESCRIPTION", "ICON", "COLOR"}
		rows := make([][]string, 0, len(views))
		for _, v := range views {
			owner := ""
			if v.Owner != nil {
				owner = v.Owner.Name
			}
			_ = owner
			rows = append(rows, []string{v.ID, v.Name, v.Description, v.Icon, v.Color})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var viewsViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр пользовательского представления по ID",
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
		view, err := c.GetCustomView(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(view)
		}

		fmt.Fprintf(out, "ID:          %s\n", view.ID)
		fmt.Fprintf(out, "Название:    %s\n", view.Name)
		if view.Description != "" {
			fmt.Fprintf(out, "Описание:    %s\n", view.Description)
		}
		if view.Icon != "" {
			fmt.Fprintf(out, "Иконка:      %s\n", view.Icon)
		}
		if view.Color != "" {
			fmt.Fprintf(out, "Цвет:        %s\n", view.Color)
		}
		if view.Owner != nil {
			fmt.Fprintf(out, "Владелец:    %s\n", view.Owner.Name)
		}
		return nil
	},
}

var viewsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать пользовательское представление",
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
		icon, _ := cmd.Flags().GetString("icon")
		color, _ := cmd.Flags().GetString("color")

		if name == "" {
			return fmt.Errorf("необходимо указать --name")
		}

		c := newLinearClient(t)
		view, err := c.CreateCustomView(name, description, icon, color)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(view)
		}
		fmt.Fprintf(out, "Представление создано: %s (%s)\n", view.Name, view.ID)
		return nil
	},
}

var viewsUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить пользовательское представление по ID",
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

		c := newLinearClient(t)
		view, err := c.UpdateCustomView(args[0], name, description)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(view)
		}
		fmt.Fprintf(out, "Представление обновлено: %s (%s)\n", view.Name, view.ID)
		return nil
	},
}

var viewsDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить пользовательское представление по ID",
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
		if err := c.DeleteCustomView(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Представление удалено: %s\n", args[0])
		return nil
	},
}

func init() {
	viewsCreateCmd.Flags().String("name", "", "Название представления (обязательно)")
	viewsCreateCmd.Flags().String("description", "", "Описание представления")
	viewsCreateCmd.Flags().String("icon", "", "Иконка представления")
	viewsCreateCmd.Flags().String("color", "", "Цвет представления")

	viewsUpdateCmd.Flags().String("name", "", "Новое название представления")
	viewsUpdateCmd.Flags().String("description", "", "Новое описание представления")

	viewsCmd.AddCommand(viewsListCmd)
	viewsCmd.AddCommand(viewsViewCmd)
	viewsCmd.AddCommand(viewsCreateCmd)
	viewsCmd.AddCommand(viewsUpdateCmd)
	viewsCmd.AddCommand(viewsDeleteCmd)

	rootCmd.AddCommand(viewsCmd)
}
