package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var initiativesCmd = &cobra.Command{
	Use:   "initiatives",
	Short: "Управление инициативами",
}

var initiativesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список инициатив",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		initiatives, err := c.ListInitiatives()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(initiatives)
		}

		if len(initiatives) == 0 {
			fmt.Fprintln(out, "Инициативы не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "STATUS", "OWNER"}
		rows := make([][]string, 0, len(initiatives))
		for _, i := range initiatives {
			ownerName := ""
			if i.Owner != nil {
				ownerName = i.Owner.Name
			}
			rows = append(rows, []string{i.ID, i.Name, i.Status, ownerName})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var initiativesViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр инициативы по ID",
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
		initiative, err := c.GetInitiative(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(initiative)
		}

		fmt.Fprintf(out, "ID:      %s\n", initiative.ID)
		fmt.Fprintf(out, "Название: %s\n", initiative.Name)
		fmt.Fprintf(out, "Статус:  %s\n", initiative.Status)
		if initiative.Owner != nil {
			fmt.Fprintf(out, "Владелец: %s\n", initiative.Owner.Name)
		}
		if initiative.Description != "" {
			fmt.Fprintf(out, "\n%s\n", initiative.Description)
		}
		return nil
	},
}

var initiativesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать инициативу",
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
		ownerID, _ := cmd.Flags().GetString("owner")

		if name == "" {
			return fmt.Errorf("необходимо указать --name")
		}

		c := newLinearClient(t)
		initiative, err := c.CreateInitiative(name, description, ownerID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(initiative)
		}
		fmt.Fprintf(out, "Инициатива создана: %s (%s)\n", initiative.Name, initiative.ID)
		return nil
	},
}

var initiativesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить инициативу по ID",
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
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			input["status"] = status
		}

		if len(input) == 0 {
			return fmt.Errorf("необходимо указать хотя бы один флаг для обновления")
		}

		c := newLinearClient(t)
		initiative, err := c.UpdateInitiative(args[0], input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(initiative)
		}
		fmt.Fprintf(out, "Инициатива обновлена: %s (%s)\n", initiative.Name, initiative.ID)
		return nil
	},
}

var initiativesArchiveCmd = &cobra.Command{
	Use:   "archive <ID>",
	Short: "Архивировать инициативу по ID",
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
		if err := c.ArchiveInitiative(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Инициатива архивирована: %s\n", args[0])
		return nil
	},
}

func init() {
	initiativesCmd.AddCommand(initiativesListCmd)
	initiativesCmd.AddCommand(initiativesViewCmd)

	initiativesCreateCmd.Flags().String("name", "", "Название инициативы (обязательно)")
	initiativesCreateCmd.Flags().String("description", "", "Описание инициативы")
	initiativesCreateCmd.Flags().String("owner", "", "ID владельца")
	initiativesCmd.AddCommand(initiativesCreateCmd)

	initiativesUpdateCmd.Flags().String("name", "", "Новое название")
	initiativesUpdateCmd.Flags().String("description", "", "Новое описание")
	initiativesUpdateCmd.Flags().String("status", "", "Новый статус")
	initiativesCmd.AddCommand(initiativesUpdateCmd)

	initiativesCmd.AddCommand(initiativesArchiveCmd)

	rootCmd.AddCommand(initiativesCmd)
}
