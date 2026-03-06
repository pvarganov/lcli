package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Управление пользователями",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список пользователей организации",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		users, err := c.ListUsers()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(users)
		}

		if len(users) == 0 {
			fmt.Fprintln(out, "Пользователи не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "DISPLAY NAME", "EMAIL"}
		rows := make([][]string, 0, len(users))
		for _, u := range users {
			rows = append(rows, []string{
				u.ID,
				format.StripControlChars(u.Name),
				format.StripControlChars(u.DisplayName),
				format.StripControlChars(u.Email),
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var usersViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр пользователя по ID",
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
		user, err := c.GetUser(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(user)
		}

		fmt.Fprintf(out, "ID:           %s\n", user.ID)
		fmt.Fprintf(out, "Name:         %s\n", user.Name)
		fmt.Fprintf(out, "Display Name: %s\n", user.DisplayName)
		fmt.Fprintf(out, "Email:        %s\n", user.Email)
		return nil
	},
}

var usersMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Просмотр текущего пользователя",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		viewer, err := c.GetViewer()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(viewer)
		}

		fmt.Fprintf(out, "ID:           %s\n", viewer.ID)
		fmt.Fprintf(out, "Name:         %s\n", viewer.Name)
		fmt.Fprintf(out, "Display Name: %s\n", viewer.DisplayName)
		fmt.Fprintf(out, "Email:        %s\n", viewer.Email)
		return nil
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersViewCmd)
	usersCmd.AddCommand(usersMeCmd)
	rootCmd.AddCommand(usersCmd)
}
