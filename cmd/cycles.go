package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var cyclesCmd = &cobra.Command{
	Use:   "cycles",
	Short: "Управление циклами (спринтами)",
}

var cyclesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список циклов команды",
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
			return fmt.Errorf("требуется --team (ID команды)")
		}

		c := newLinearClient(t)
		cycles, err := c.ListCycles(teamID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(cycles)
		}

		if len(cycles) == 0 {
			fmt.Fprintln(out, "Циклы не найдены.")
			return nil
		}

		headers := []string{"ID", "NUMBER", "NAME", "STARTS AT", "ENDS AT"}
		rows := make([][]string, 0, len(cycles))
		for _, cy := range cycles {
			rows = append(rows, []string{
				cy.ID,
				fmt.Sprintf("%d", cy.Number),
				format.StripControlChars(cy.Name),
				cy.StartsAt.Format("2006-01-02"),
				cy.EndsAt.Format("2006-01-02"),
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var cyclesViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр цикла по ID",
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
		cycle, err := c.GetCycle(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(cycle)
		}

		fmt.Fprintf(out, "ID:      %s\n", cycle.ID)
		fmt.Fprintf(out, "Number:  %d\n", cycle.Number)
		fmt.Fprintf(out, "Name:    %s\n", format.StripControlChars(cycle.Name))
		fmt.Fprintf(out, "Team:    %s (%s)\n", format.StripControlChars(cycle.Team.Name), cycle.Team.Key)
		fmt.Fprintf(out, "Starts:  %s\n", cycle.StartsAt.Format("2006-01-02"))
		fmt.Fprintf(out, "Ends:    %s\n", cycle.EndsAt.Format("2006-01-02"))
		if cycle.CompletedAt != nil {
			fmt.Fprintf(out, "Completed: %s\n", cycle.CompletedAt.Format("2006-01-02"))
		}
		return nil
	},
}

var cyclesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новый цикл",
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
			return fmt.Errorf("требуется --team (ID команды)")
		}
		startsAt, _ := cmd.Flags().GetString("starts-at")
		if startsAt == "" {
			return fmt.Errorf("требуется --starts-at (дата начала, формат: 2006-01-02T00:00:00Z)")
		}
		endsAt, _ := cmd.Flags().GetString("ends-at")
		if endsAt == "" {
			return fmt.Errorf("требуется --ends-at (дата окончания, формат: 2006-01-02T00:00:00Z)")
		}
		name, _ := cmd.Flags().GetString("name")

		c := newLinearClient(t)
		cycle, err := c.CreateCycle(teamID, name, startsAt, endsAt)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(cycle)
		}

		fmt.Fprintf(out, "Цикл создан: %s (ID: %s)\n", format.StripControlChars(cycle.Name), cycle.ID)
		return nil
	},
}

var cyclesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить цикл",
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
		startsAt, _ := cmd.Flags().GetString("starts-at")
		endsAt, _ := cmd.Flags().GetString("ends-at")

		c := newLinearClient(t)
		cycle, err := c.UpdateCycle(args[0], name, startsAt, endsAt)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(cycle)
		}

		fmt.Fprintf(out, "Цикл обновлён: %s (ID: %s)\n", format.StripControlChars(cycle.Name), cycle.ID)
		return nil
	},
}

var cyclesArchiveCmd = &cobra.Command{
	Use:   "archive <ID>",
	Short: "Архивировать цикл",
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
		if err := c.ArchiveCycle(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Цикл %s архивирован.\n", args[0])
		return nil
	},
}

func init() {
	cyclesListCmd.Flags().String("team", "", "ID команды (обязательно)")

	cyclesCreateCmd.Flags().String("team", "", "ID команды (обязательно)")
	cyclesCreateCmd.Flags().String("name", "", "Название цикла")
	cyclesCreateCmd.Flags().String("starts-at", "", "Дата начала (формат: 2006-01-02T00:00:00Z)")
	cyclesCreateCmd.Flags().String("ends-at", "", "Дата окончания (формат: 2006-01-02T00:00:00Z)")

	cyclesUpdateCmd.Flags().String("name", "", "Новое название цикла")
	cyclesUpdateCmd.Flags().String("starts-at", "", "Новая дата начала")
	cyclesUpdateCmd.Flags().String("ends-at", "", "Новая дата окончания")

	cyclesCmd.AddCommand(cyclesListCmd)
	cyclesCmd.AddCommand(cyclesViewCmd)
	cyclesCmd.AddCommand(cyclesCreateCmd)
	cyclesCmd.AddCommand(cyclesUpdateCmd)
	cyclesCmd.AddCommand(cyclesArchiveCmd)
	rootCmd.AddCommand(cyclesCmd)
}
