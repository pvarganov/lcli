package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Управление вебхуками",
}

var webhooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список вебхуков",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		webhooks, err := c.ListWebhooks()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(webhooks)
		}

		if len(webhooks) == 0 {
			fmt.Fprintln(out, "Вебхуки не найдены.")
			return nil
		}

		headers := []string{"ID", "URL", "ENABLED", "RESOURCE TYPES", "TEAM"}
		rows := make([][]string, 0, len(webhooks))
		for _, wh := range webhooks {
			teamKey := ""
			if wh.Team != nil {
				teamKey = wh.Team.Key
			}
			enabled := "no"
			if wh.Enabled {
				enabled = "yes"
			}
			rows = append(rows, []string{
				wh.ID,
				format.StripControlChars(wh.URL),
				enabled,
				strings.Join(wh.ResourceTypes, ","),
				teamKey,
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var webhooksViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр вебхука по ID",
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
		wh, err := c.GetWebhook(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(wh)
		}

		enabled := "нет"
		if wh.Enabled {
			enabled = "да"
		}
		teamKey := ""
		if wh.Team != nil {
			teamKey = wh.Team.Key
		}
		fmt.Fprintf(out, "ID:             %s\n", wh.ID)
		fmt.Fprintf(out, "URL:            %s\n", wh.URL)
		fmt.Fprintf(out, "Enabled:        %s\n", enabled)
		fmt.Fprintf(out, "Secret:         %s\n", wh.Secret)
		fmt.Fprintf(out, "ResourceTypes:  %s\n", strings.Join(wh.ResourceTypes, ", "))
		fmt.Fprintf(out, "Team:           %s\n", teamKey)
		return nil
	},
}

var webhooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать вебхук",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		url, _ := cmd.Flags().GetString("url")
		teamID, _ := cmd.Flags().GetString("team-id")
		secret, _ := cmd.Flags().GetString("secret")
		resourceTypesStr, _ := cmd.Flags().GetString("resource-types")
		enabled, _ := cmd.Flags().GetBool("enabled")

		if url == "" {
			return fmt.Errorf("флаг --url обязателен")
		}
		if teamID == "" {
			return fmt.Errorf("флаг --team-id обязателен")
		}

		var resourceTypes []string
		if resourceTypesStr != "" {
			resourceTypes = strings.Split(resourceTypesStr, ",")
		}

		c := newLinearClient(t)
		wh, err := c.CreateWebhook(client.WebhookCreateInput{
			URL:           url,
			TeamID:        teamID,
			Enabled:       enabled,
			Secret:        secret,
			ResourceTypes: resourceTypes,
		})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(wh)
		}
		fmt.Fprintf(out, "Вебхук создан: %s (%s)\n", wh.ID, wh.URL)
		return nil
	},
}

var webhooksUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить вебхук",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		inp := map[string]any{}
		if url, _ := cmd.Flags().GetString("url"); url != "" {
			inp["url"] = url
		}
		if cmd.Flags().Changed("enabled") {
			enabled, _ := cmd.Flags().GetBool("enabled")
			inp["enabled"] = enabled
		}
		if secret, _ := cmd.Flags().GetString("secret"); secret != "" {
			inp["secret"] = secret
		}
		if resourceTypesStr, _ := cmd.Flags().GetString("resource-types"); resourceTypesStr != "" {
			inp["resourceTypes"] = strings.Split(resourceTypesStr, ",")
		}

		if len(inp) == 0 {
			return fmt.Errorf("укажите хотя бы один флаг для обновления")
		}

		c := newLinearClient(t)
		wh, err := c.UpdateWebhook(args[0], inp)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(wh)
		}
		fmt.Fprintf(out, "Вебхук %s обновлён.\n", wh.ID)
		return nil
	},
}

var webhooksDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить вебхук",
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
		if err := c.DeleteWebhook(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Вебхук %s удалён.\n", args[0])
		return nil
	},
}

var webhooksRotateSecretCmd = &cobra.Command{
	Use:   "rotate-secret <ID>",
	Short: "Обновить секрет вебхука",
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
		secret, err := c.RotateWebhookSecret(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(map[string]string{"id": args[0], "secret": secret})
		}
		fmt.Fprintf(out, "Секрет вебхука %s обновлён. Новый секрет: %s\n", args[0], secret)
		return nil
	},
}

func init() {
	webhooksCreateCmd.Flags().String("url", "", "URL вебхука (обязательно)")
	webhooksCreateCmd.Flags().String("team-id", "", "ID команды (обязательно)")
	webhooksCreateCmd.Flags().String("secret", "", "Секрет вебхука")
	webhooksCreateCmd.Flags().String("resource-types", "", "Типы ресурсов через запятую (например: Issue,Comment)")
	webhooksCreateCmd.Flags().Bool("enabled", true, "Включить вебхук")

	webhooksUpdateCmd.Flags().String("url", "", "Новый URL вебхука")
	webhooksUpdateCmd.Flags().String("secret", "", "Новый секрет вебхука")
	webhooksUpdateCmd.Flags().String("resource-types", "", "Новые типы ресурсов через запятую")
	webhooksUpdateCmd.Flags().Bool("enabled", true, "Включить/отключить вебхук")

	webhooksCmd.AddCommand(webhooksListCmd)
	webhooksCmd.AddCommand(webhooksViewCmd)
	webhooksCmd.AddCommand(webhooksCreateCmd)
	webhooksCmd.AddCommand(webhooksUpdateCmd)
	webhooksCmd.AddCommand(webhooksDeleteCmd)
	webhooksCmd.AddCommand(webhooksRotateSecretCmd)
	rootCmd.AddCommand(webhooksCmd)
}
