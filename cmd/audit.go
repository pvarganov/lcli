package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Аудит-лог организации",
}

var auditListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список записей аудит-лога",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		limit, _ := cmd.Flags().GetInt("limit")
		afterCursor, _ := cmd.Flags().GetString("after")
		typeFilter, _ := cmd.Flags().GetString("type")

		c := newLinearClient(t)
		entries, _, err := c.ListAuditEntries(client.AuditEntryFilter{
			Limit: limit,
			After: afterCursor,
			Type:  typeFilter,
		})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(entries)
		}

		if len(entries) == 0 {
			fmt.Fprintln(out, "Записи аудит-лога не найдены.")
			return nil
		}

		headers := []string{"ID", "TYPE", "ACTOR_ID", "IP", "COUNTRY", "CREATED_AT"}
		rows := make([][]string, 0, len(entries))
		for _, e := range entries {
			rows = append(rows, []string{e.ID, e.Type, e.ActorID, e.IP, e.Country, e.CreatedAt})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var auditTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "Список типов аудит-лога",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		types, err := c.ListAuditEntryTypes()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(types)
		}

		if len(types) == 0 {
			fmt.Fprintln(out, "Типы аудит-лога не найдены.")
			return nil
		}

		headers := []string{"#", "TYPE", "DESCRIPTION"}
		rows := make([][]string, 0, len(types))
		for i, typ := range types {
			rows = append(rows, []string{strconv.Itoa(i + 1), typ.Type, typ.Description})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

func init() {
	auditListCmd.Flags().Int("limit", 50, "Максимальное количество записей")
	auditListCmd.Flags().String("after", "", "Курсор для пагинации")
	auditListCmd.Flags().String("type", "", "Фильтр по типу события")

	auditCmd.AddCommand(auditListCmd)
	auditCmd.AddCommand(auditTypesCmd)

	rootCmd.AddCommand(auditCmd)
}
