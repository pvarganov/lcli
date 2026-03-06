package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
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

// --- updates subcommands ---

var initiativeUpdatesCmd = &cobra.Command{
	Use:   "updates",
	Short: "Управление обновлениями инициативы",
}

var initiativeUpdatesListCmd = &cobra.Command{
	Use:   "list <INITIATIVE-ID>",
	Short: "Список обновлений инициативы",
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
		updates, err := c.ListInitiativeUpdates(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(updates)
		}

		if len(updates) == 0 {
			fmt.Fprintln(out, "Обновления не найдены.")
			return nil
		}

		headers := []string{"ID", "HEALTH", "AUTHOR", "CREATED AT", "BODY"}
		rows := make([][]string, 0, len(updates))
		for _, u := range updates {
			author := ""
			if u.User != nil {
				author = format.StripControlChars(u.User.DisplayName)
			}
			body := format.StripControlChars(u.Body)
			if len(body) > 60 {
				body = body[:57] + "..."
			}
			rows = append(rows, []string{
				u.ID,
				u.Health,
				author,
				u.CreatedAt.Format("2006-01-02"),
				body,
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var initiativeUpdatesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать обновление инициативы",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		initiativeID, _ := cmd.Flags().GetString("initiative-id")
		body, _ := cmd.Flags().GetString("body")
		health, _ := cmd.Flags().GetString("health")

		if initiativeID == "" {
			return fmt.Errorf("требуется --initiative-id")
		}

		c := newLinearClient(t)
		update, err := c.CreateInitiativeUpdate(client.CreateInitiativeUpdateInput{
			InitiativeID: initiativeID,
			Body:         body,
			Health:       health,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Обновление создано: %s\n", update.ID)
		return nil
	},
}

var initiativeUpdatesUpdateCmd = &cobra.Command{
	Use:   "update <UPDATE-ID>",
	Short: "Изменить обновление инициативы",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		body, _ := cmd.Flags().GetString("body")
		health, _ := cmd.Flags().GetString("health")

		c := newLinearClient(t)
		update, err := c.UpdateInitiativeUpdate(args[0], client.UpdateInitiativeUpdateInput{
			Body:   body,
			Health: health,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Обновление изменено: %s\n", update.ID)
		return nil
	},
}

var initiativeUpdatesDeleteCmd = &cobra.Command{
	Use:   "delete <UPDATE-ID>",
	Short: "Архивировать обновление инициативы",
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
		if err := c.ArchiveInitiativeUpdate(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Обновление удалено: %s\n", args[0])
		return nil
	},
}

// --- link-project / unlink-project ---

var initiativesLinkProjectCmd = &cobra.Command{
	Use:   "link-project",
	Short: "Связать инициативу с проектом",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		initiativeID, _ := cmd.Flags().GetString("initiative")
		projectID, _ := cmd.Flags().GetString("project")

		if initiativeID == "" {
			return fmt.Errorf("требуется --initiative")
		}
		if projectID == "" {
			return fmt.Errorf("требуется --project")
		}

		c := newLinearClient(t)
		rel, err := c.CreateInitiativeToProject(initiativeID, projectID)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Связь создана: %s (%s → %s)\n", rel.ID, rel.Initiative.Name, rel.Project.Name)
		return nil
	},
}

var initiativesUnlinkProjectCmd = &cobra.Command{
	Use:   "unlink-project <RELATION-ID>",
	Short: "Удалить связь инициативы с проектом",
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
		if err := c.DeleteInitiativeToProject(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Связь удалена: %s\n", args[0])
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

	// updates subcommands
	initiativeUpdatesCreateCmd.Flags().String("initiative-id", "", "ID инициативы (обязательно)")
	initiativeUpdatesCreateCmd.Flags().String("body", "", "Текст обновления")
	initiativeUpdatesCreateCmd.Flags().String("health", "", "Статус здоровья (onTrack, atRisk, offTrack)")

	initiativeUpdatesUpdateCmd.Flags().String("body", "", "Новый текст обновления")
	initiativeUpdatesUpdateCmd.Flags().String("health", "", "Новый статус здоровья")

	initiativeUpdatesCmd.AddCommand(initiativeUpdatesListCmd)
	initiativeUpdatesCmd.AddCommand(initiativeUpdatesCreateCmd)
	initiativeUpdatesCmd.AddCommand(initiativeUpdatesUpdateCmd)
	initiativeUpdatesCmd.AddCommand(initiativeUpdatesDeleteCmd)
	initiativesCmd.AddCommand(initiativeUpdatesCmd)

	// link/unlink project
	initiativesLinkProjectCmd.Flags().String("initiative", "", "ID инициативы (обязательно)")
	initiativesLinkProjectCmd.Flags().String("project", "", "ID проекта (обязательно)")
	_ = initiativesLinkProjectCmd.MarkFlagRequired("initiative")
	_ = initiativesLinkProjectCmd.MarkFlagRequired("project")
	initiativesCmd.AddCommand(initiativesLinkProjectCmd)
	initiativesCmd.AddCommand(initiativesUnlinkProjectCmd)

	rootCmd.AddCommand(initiativesCmd)
}
