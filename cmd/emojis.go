package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var emojisCmd = &cobra.Command{
	Use:   "emojis",
	Short: "Управление эмодзи организации",
}

var emojisListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список эмодзи организации",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		emojis, err := c.ListEmojis()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(emojis)
		}

		if len(emojis) == 0 {
			fmt.Fprintln(out, "Эмодзи не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "URL"}
		rows := make([][]string, 0, len(emojis))
		for _, e := range emojis {
			rows = append(rows, []string{e.ID, e.Name, e.URL})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var emojisCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать кастомный эмодзи",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		url, _ := cmd.Flags().GetString("url")

		if name == "" {
			return fmt.Errorf("необходимо указать --name")
		}
		if url == "" {
			return fmt.Errorf("необходимо указать --url")
		}

		c := newLinearClient(t)
		emoji, err := c.CreateEmoji(name, url)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(emoji)
		}
		fmt.Fprintf(out, "Эмодзи создан: %s (id: %s)\n", emoji.Name, emoji.ID)
		return nil
	},
}

var emojisDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить кастомный эмодзи по ID",
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
		if err := c.DeleteEmoji(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Эмодзи удалён: %s\n", args[0])
		return nil
	},
}

func init() {
	emojisCreateCmd.Flags().String("name", "", "Название эмодзи (обязательно)")
	emojisCreateCmd.Flags().String("url", "", "URL изображения эмодзи (обязательно)")

	emojisCmd.AddCommand(emojisListCmd)
	emojisCmd.AddCommand(emojisCreateCmd)
	emojisCmd.AddCommand(emojisDeleteCmd)

	rootCmd.AddCommand(emojisCmd)
}
