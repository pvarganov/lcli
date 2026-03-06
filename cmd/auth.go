package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pavelvarganov/lcli/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Управление аутентификацией",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Сохранить Linear API токен в конфиг-файл",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprint(cmd.OutOrStdout(), "Введите Linear API токен: ")
		reader := bufio.NewReader(os.Stdin)
		t, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read token: %w", err)
		}
		t = strings.TrimSpace(t)
		if t == "" {
			return fmt.Errorf("токен не может быть пустым")
		}
		if err := config.SaveToken(t); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Токен сохранён успешно.")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Показать текущий токен аутентификации",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Приоритет: флаг --token → LINEAR_API_KEY → файл
		t := GetToken()
		if t == "" {
			var err error
			t, err = config.LoadToken()
			if err != nil {
				return err
			}
		}
		if t == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "Токен не настроен. Запустите `lcli auth login`.")
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Токен: %s\n", config.MaskToken(t))
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
