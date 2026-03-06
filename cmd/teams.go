package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Управление командами",
}

var teamsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список команд",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		teams, err := c.ListTeams()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(teams)
		}

		if len(teams) == 0 {
			fmt.Fprintln(out, "Команды не найдены.")
			return nil
		}

		headers := []string{"ID", "KEY", "NAME"}
		rows := make([][]string, 0, len(teams))
		for _, team := range teams {
			rows = append(rows, []string{team.ID, format.StripControlChars(team.Key), format.StripControlChars(team.Name)})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var teamsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать команду",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")
		if name == "" {
			return fmt.Errorf("требуется --name (название команды)")
		}
		if key == "" {
			return fmt.Errorf("требуется --key (ключ команды)")
		}

		c := newLinearClient(t)
		team, err := c.CreateTeam(name, key)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(team)
		}

		fmt.Fprintf(out, "Команда создана: [%s] %s (ID: %s)\n", team.Key, team.Name, team.ID)
		return nil
	},
}

var teamsUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить команду",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")

		c := newLinearClient(t)
		team, err := c.UpdateTeam(args[0], client.UpdateTeamInput{Name: name, Key: key})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(team)
		}

		fmt.Fprintf(out, "Команда обновлена: [%s] %s\n", team.Key, team.Name)
		return nil
	},
}

var teamsDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить команду",
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
		if err := c.DeleteTeam(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Команда %s удалена\n", args[0])
		return nil
	},
}

func init() {
	teamsCreateCmd.Flags().String("name", "", "Название команды")
	teamsCreateCmd.Flags().String("key", "", "Ключ команды (например ENG)")

	teamsUpdateCmd.Flags().String("name", "", "Новое название команды")
	teamsUpdateCmd.Flags().String("key", "", "Новый ключ команды")

	teamsCmd.AddCommand(teamsListCmd)
	teamsCmd.AddCommand(teamsCreateCmd)
	teamsCmd.AddCommand(teamsUpdateCmd)
	teamsCmd.AddCommand(teamsDeleteCmd)
	rootCmd.AddCommand(teamsCmd)
}
