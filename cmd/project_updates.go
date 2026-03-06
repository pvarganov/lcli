package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var projectUpdatesCmd = &cobra.Command{
	Use:   "updates",
	Short: "Управление обновлениями проекта",
}

var projectUpdatesListCmd = &cobra.Command{
	Use:   "list <PROJECT-ID>",
	Short: "Список обновлений проекта",
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
		updates, err := c.ListProjectUpdates(args[0])
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

var projectUpdatesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать обновление проекта",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		projectID, _ := cmd.Flags().GetString("project-id")
		body, _ := cmd.Flags().GetString("body")
		health, _ := cmd.Flags().GetString("health")

		if projectID == "" {
			return fmt.Errorf("требуется --project-id")
		}
		if body == "" {
			return fmt.Errorf("требуется --body")
		}

		c := newLinearClient(t)
		update, err := c.CreateProjectUpdate(client.CreateProjectUpdateInput{
			ProjectID: projectID,
			Body:      body,
			Health:    health,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Обновление создано: %s\n", update.ID)
		return nil
	},
}

var projectUpdatesUpdateCmd = &cobra.Command{
	Use:   "update <UPDATE-ID>",
	Short: "Обновить запись журнала проекта",
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
		update, err := c.UpdateProjectUpdate(args[0], client.UpdateProjectUpdateInput{
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

var projectUpdatesDeleteCmd = &cobra.Command{
	Use:   "delete <UPDATE-ID>",
	Short: "Удалить (архивировать) обновление проекта",
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
		if err := c.ArchiveProjectUpdate(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Обновление удалено: %s\n", args[0])
		return nil
	},
}

func init() {
	projectUpdatesCreateCmd.Flags().String("project-id", "", "ID проекта (обязательно)")
	projectUpdatesCreateCmd.Flags().String("body", "", "Текст обновления (обязательно)")
	projectUpdatesCreateCmd.Flags().String("health", "", "Статус здоровья проекта (onTrack, atRisk, offTrack)")

	projectUpdatesUpdateCmd.Flags().String("body", "", "Новый текст обновления")
	projectUpdatesUpdateCmd.Flags().String("health", "", "Новый статус здоровья")

	projectUpdatesCmd.AddCommand(projectUpdatesListCmd)
	projectUpdatesCmd.AddCommand(projectUpdatesCreateCmd)
	projectUpdatesCmd.AddCommand(projectUpdatesUpdateCmd)
	projectUpdatesCmd.AddCommand(projectUpdatesDeleteCmd)
	projectsCmd.AddCommand(projectUpdatesCmd)
}
