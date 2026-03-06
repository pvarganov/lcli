package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Управление организацией",
}

var orgViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Просмотр информации об организации",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		org, err := c.GetOrganization()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(org)
		}

		fmt.Fprintf(out, "ID:          %s\n", org.ID)
		fmt.Fprintf(out, "Название:    %s\n", org.Name)
		fmt.Fprintf(out, "URL Key:     %s\n", org.URLKey)
		if org.LogoURL != "" {
			fmt.Fprintf(out, "Логотип:     %s\n", org.LogoURL)
		}
		if !org.CreatedAt.IsZero() {
			fmt.Fprintf(out, "Создана:     %s\n", org.CreatedAt.Format("2006-01-02"))
		}
		return nil
	},
}

var orgInvitesCmd = &cobra.Command{
	Use:   "invites",
	Short: "Управление приглашениями в организацию",
}

var orgInvitesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список приглашений в организацию",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		invites, err := c.ListOrganizationInvites()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(invites)
		}

		if len(invites) == 0 {
			fmt.Fprintln(out, "Приглашения не найдены.")
			return nil
		}

		headers := []string{"ID", "EMAIL", "ROLE"}
		rows := make([][]string, 0, len(invites))
		for _, inv := range invites {
			rows = append(rows, []string{inv.ID, inv.Email, inv.Role})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var orgInvitesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать приглашение в организацию",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		email, _ := cmd.Flags().GetString("email")
		role, _ := cmd.Flags().GetString("role")

		if email == "" {
			return fmt.Errorf("необходимо указать --email")
		}

		c := newLinearClient(t)
		inv, err := c.CreateOrganizationInvite(email, role)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(inv)
		}
		fmt.Fprintf(out, "Приглашение создано: %s (%s)\n", inv.Email, inv.ID)
		return nil
	},
}

var orgInvitesDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить приглашение по ID",
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
		if err := c.DeleteOrganizationInvite(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Приглашение удалено: %s\n", args[0])
		return nil
	},
}

var orgInvitesResendCmd = &cobra.Command{
	Use:   "resend <ID>",
	Short: "Повторно отправить приглашение по ID",
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
		if err := c.ResendOrganizationInvite(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Приглашение повторно отправлено: %s\n", args[0])
		return nil
	},
}

func init() {
	orgCmd.AddCommand(orgViewCmd)

	orgInvitesCreateCmd.Flags().String("email", "", "Email приглашаемого пользователя (обязательно)")
	orgInvitesCreateCmd.Flags().String("role", "", "Роль пользователя (например: member, admin)")

	orgInvitesCmd.AddCommand(orgInvitesListCmd)
	orgInvitesCmd.AddCommand(orgInvitesCreateCmd)
	orgInvitesCmd.AddCommand(orgInvitesDeleteCmd)
	orgInvitesCmd.AddCommand(orgInvitesResendCmd)
	orgCmd.AddCommand(orgInvitesCmd)

	rootCmd.AddCommand(orgCmd)
}
