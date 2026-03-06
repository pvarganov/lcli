package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var integrationsCmd = &cobra.Command{
	Use:   "integrations",
	Short: "Управление интеграциями",
}

var integrationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список интеграций",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		integrations, err := c.ListIntegrations()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(integrations)
		}

		if len(integrations) == 0 {
			fmt.Fprintln(out, "Интеграции не найдены.")
			return nil
		}

		headers := []string{"ID", "SERVICE", "TEAM", "CREATED_AT"}
		rows := make([][]string, 0, len(integrations))
		for _, intg := range integrations {
			team := ""
			if intg.Team != nil {
				team = intg.Team.Key
			}
			rows = append(rows, []string{intg.ID, intg.Service, team, intg.CreatedAt})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var integrationsDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить интеграцию",
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
		if err := c.DeleteIntegration(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Интеграция удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	integrationsCmd.AddCommand(integrationsListCmd)
	integrationsCmd.AddCommand(integrationsDeleteCmd)

	rootCmd.AddCommand(integrationsCmd)
}
