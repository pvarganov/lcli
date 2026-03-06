package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var teamMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Управление членами команды",
}

var teamMembersListCmd = &cobra.Command{
	Use:   "list <TEAM-ID>",
	Short: "Список членов команды",
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
		members, err := c.ListTeamMembers(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(members)
		}

		if len(members) == 0 {
			fmt.Fprintln(out, "Члены команды не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "EMAIL", "ROLE"}
		rows := make([][]string, 0, len(members))
		for _, m := range members {
			rows = append(rows, []string{
				m.ID,
				format.StripControlChars(m.User.DisplayName),
				format.StripControlChars(m.User.Email),
				format.StripControlChars(m.Role),
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var teamMembersAddCmd = &cobra.Command{
	Use:   "add <TEAM-ID> <USER-ID>",
	Short: "Добавить пользователя в команду",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		m, err := c.CreateTeamMembership(args[0], args[1])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(m)
		}

		fmt.Fprintf(out, "Пользователь %s добавлен в команду %s\n", m.User.DisplayName, m.Team.Name)
		return nil
	},
}

var teamMembersRemoveCmd = &cobra.Command{
	Use:   "remove <MEMBERSHIP-ID>",
	Short: "Удалить пользователя из команды",
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
		if err := c.DeleteTeamMembership(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Членство %s удалено\n", args[0])
		return nil
	},
}

func init() {
	teamMembersCmd.AddCommand(teamMembersListCmd)
	teamMembersCmd.AddCommand(teamMembersAddCmd)
	teamMembersCmd.AddCommand(teamMembersRemoveCmd)
	teamsCmd.AddCommand(teamMembersCmd)
}
