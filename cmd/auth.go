package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pavelvarganov/lcli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
		var t string
		if term.IsTerminal(int(os.Stdin.Fd())) {
			raw, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("read token: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			t = strings.TrimSpace(string(raw))
		} else {
			reader := bufio.NewReader(os.Stdin)
			raw, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read token: %w", err)
			}
			t = strings.TrimSpace(raw)
		}
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
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
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
