package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var rateLimitCmd = &cobra.Command{
	Use:   "rate-limit",
	Short: "Статус лимита запросов API",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		status, err := c.GetRateLimitStatus()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(status)
		}

		fmt.Fprintf(out, "Identifier: %s\nKind: %s\n\n", status.Identifier, status.Kind)
		if len(status.Limits) > 0 {
			headers := []string{"TYPE", "ALLOWED", "REQUESTED", "REMAINING", "RESET"}
			rows := make([][]string, 0, len(status.Limits))
			for _, l := range status.Limits {
				rows = append(rows, []string{
					l.Type,
					strconv.FormatFloat(l.AllowedAmount, 'f', 0, 64),
					strconv.FormatFloat(l.RequestedAmount, 'f', 0, 64),
					strconv.FormatFloat(l.RemainingAmount, 'f', 0, 64),
					l.Reset,
				})
			}
			format.TableWriter(out, headers, rows)
		}
		return nil
	},
}

var timeSchedulesCmd = &cobra.Command{
	Use:   "time-schedules",
	Short: "Управление расписаниями",
}

var timeSchedulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список расписаний",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		schedules, err := c.ListTimeSchedules()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(schedules)
		}

		if len(schedules) == 0 {
			fmt.Fprintln(out, "Расписания не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "CREATED_AT"}
		rows := make([][]string, 0, len(schedules))
		for _, s := range schedules {
			rows = append(rows, []string{s.ID, s.Name, s.CreatedAt})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var timeSchedulesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать расписание",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")

		if name == "" {
			return fmt.Errorf("--name обязателен")
		}

		c := newLinearClient(t)
		ts, err := c.CreateTimeSchedule(name)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(ts)
		}

		fmt.Fprintf(out, "Расписание создано: %s (%s)\n", ts.Name, ts.ID)
		return nil
	},
}

var timeSchedulesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить расписание",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		id := args[0]
		name, _ := cmd.Flags().GetString("name")

		c := newLinearClient(t)
		ts, err := c.UpdateTimeSchedule(id, name)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(ts)
		}

		fmt.Fprintf(out, "Расписание обновлено: %s (%s)\n", ts.Name, ts.ID)
		return nil
	},
}

var timeSchedulesDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить расписание",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		id := args[0]
		c := newLinearClient(t)
		if err := c.DeleteTimeSchedule(id); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Расписание удалено: %s\n", id)
		return nil
	},
}

var triageResponsibilitiesCmd = &cobra.Command{
	Use:   "triage-responsibilities",
	Short: "Управление ответственностью за триаж",
}

var triageResponsibilitiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список ответственных за триаж",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		teamID, _ := cmd.Flags().GetString("team")
		if teamID == "" {
			return fmt.Errorf("--team обязателен")
		}

		c := newLinearClient(t)
		responsibilities, err := c.ListTriageResponsibilities(teamID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(responsibilities)
		}

		if len(responsibilities) == 0 {
			fmt.Fprintln(out, "Ответственные за триаж не найдены.")
			return nil
		}

		headers := []string{"ID", "ACTION", "TEAM", "USER", "CREATED_AT"}
		rows := make([][]string, 0, len(responsibilities))
		for _, r := range responsibilities {
			userName := ""
			if r.CurrentUser != nil {
				userName = r.CurrentUser.Name
			}
			rows = append(rows, []string{r.ID, r.Action, r.Team.Key, userName, r.CreatedAt})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

func init() {
	timeSchedulesCreateCmd.Flags().String("name", "", "Название расписания")

	timeSchedulesUpdateCmd.Flags().String("name", "", "Новое название расписания")

	triageResponsibilitiesListCmd.Flags().String("team", "", "ID команды")

	timeSchedulesCmd.AddCommand(timeSchedulesListCmd)
	timeSchedulesCmd.AddCommand(timeSchedulesCreateCmd)
	timeSchedulesCmd.AddCommand(timeSchedulesUpdateCmd)
	timeSchedulesCmd.AddCommand(timeSchedulesDeleteCmd)

	triageResponsibilitiesCmd.AddCommand(triageResponsibilitiesListCmd)

	rootCmd.AddCommand(rateLimitCmd)
	rootCmd.AddCommand(timeSchedulesCmd)
	rootCmd.AddCommand(triageResponsibilitiesCmd)
}

