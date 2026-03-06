package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var gitAutomationCmd = &cobra.Command{
	Use:   "git-automation",
	Short: "Управление git автоматизацией",
}

var gitAutomationStatesCmd = &cobra.Command{
	Use:   "states",
	Short: "Управление правилами git автоматизации",
}

var gitAutomationStatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список правил git автоматизации команды",
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
		states, err := c.ListGitAutomationStates(teamID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(states)
		}

		if len(states) == 0 {
			fmt.Fprintln(out, "Правила git автоматизации не найдены.")
			return nil
		}

		headers := []string{"ID", "EVENT", "STATE", "TARGET_BRANCH", "TEAM"}
		rows := make([][]string, 0, len(states))
		for _, s := range states {
			stateName := ""
			if s.State != nil {
				stateName = format.StripControlChars(s.State.Name)
			}
			targetBranch := ""
			if s.TargetBranch != nil {
				targetBranch = s.TargetBranch.BranchPattern
			} else if s.BranchPattern != "" {
				targetBranch = s.BranchPattern
			}
			teamKey := ""
			if s.Team != nil {
				teamKey = s.Team.Key
			}
			rows = append(rows, []string{s.ID, s.Event, stateName, targetBranch, teamKey})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var gitAutomationStatesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать правило git автоматизации",
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
		event, _ := cmd.Flags().GetString("event")
		if event == "" {
			return fmt.Errorf("требуется --event (например: branchCreated, branchMerged)")
		}
		stateID, _ := cmd.Flags().GetString("state")
		targetBranchID, _ := cmd.Flags().GetString("target-branch")

		c := newLinearClient(t)
		state, err := c.CreateGitAutomationState(teamID, event, stateID, targetBranchID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(state)
		}

		fmt.Fprintf(out, "Правило создано: %s (событие: %s)\n", state.ID, state.Event)
		return nil
	},
}

var gitAutomationStatesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить правило git автоматизации",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		event, _ := cmd.Flags().GetString("event")
		stateID, _ := cmd.Flags().GetString("state")
		targetBranchID, _ := cmd.Flags().GetString("target-branch")

		c := newLinearClient(t)
		state, err := c.UpdateGitAutomationState(args[0], event, stateID, targetBranchID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(state)
		}

		fmt.Fprintf(out, "Правило обновлено: %s (событие: %s)\n", state.ID, state.Event)
		return nil
	},
}

var gitAutomationStatesDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить правило git автоматизации",
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
		if err := c.DeleteGitAutomationState(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Правило %s удалено.\n", args[0])
		return nil
	},
}

var gitAutomationBranchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "Управление целевыми ветками git автоматизации",
}

var gitAutomationBranchesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать целевую ветку для git автоматизации",
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
		pattern, _ := cmd.Flags().GetString("pattern")
		if pattern == "" {
			return fmt.Errorf("требуется --pattern (паттерн ветки)")
		}
		isRegex, _ := cmd.Flags().GetBool("regex")

		c := newLinearClient(t)
		tb, err := c.CreateGitAutomationTargetBranch(teamID, pattern, isRegex)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(tb)
		}

		fmt.Fprintf(out, "Целевая ветка создана: %s (паттерн: %s)\n", tb.ID, tb.BranchPattern)
		return nil
	},
}

var gitAutomationBranchesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить целевую ветку для git автоматизации",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		pattern, _ := cmd.Flags().GetString("pattern")
		isRegex, _ := cmd.Flags().GetBool("regex")

		c := newLinearClient(t)
		tb, err := c.UpdateGitAutomationTargetBranch(args[0], pattern, isRegex)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(tb)
		}

		fmt.Fprintf(out, "Целевая ветка обновлена: %s (паттерн: %s)\n", tb.ID, tb.BranchPattern)
		return nil
	},
}

var gitAutomationBranchesDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить целевую ветку для git автоматизации",
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
		if err := c.DeleteGitAutomationTargetBranch(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Целевая ветка %s удалена.\n", args[0])
		return nil
	},
}

func init() {
	gitAutomationStatesListCmd.Flags().String("team", "", "ID команды (обязательно)")

	gitAutomationStatesCreateCmd.Flags().String("team", "", "ID команды (обязательно)")
	gitAutomationStatesCreateCmd.Flags().String("event", "", "Событие (например: branchCreated, branchMerged, branchDeleted)")
	gitAutomationStatesCreateCmd.Flags().String("state", "", "ID статуса workflow")
	gitAutomationStatesCreateCmd.Flags().String("target-branch", "", "ID целевой ветки")

	gitAutomationStatesUpdateCmd.Flags().String("event", "", "Событие")
	gitAutomationStatesUpdateCmd.Flags().String("state", "", "ID статуса workflow")
	gitAutomationStatesUpdateCmd.Flags().String("target-branch", "", "ID целевой ветки")

	gitAutomationBranchesCreateCmd.Flags().String("team", "", "ID команды (обязательно)")
	gitAutomationBranchesCreateCmd.Flags().String("pattern", "", "Паттерн ветки (обязательно)")
	gitAutomationBranchesCreateCmd.Flags().Bool("regex", false, "Использовать regex паттерн")

	gitAutomationBranchesUpdateCmd.Flags().String("pattern", "", "Новый паттерн ветки")
	gitAutomationBranchesUpdateCmd.Flags().Bool("regex", false, "Использовать regex паттерн")

	gitAutomationStatesCmd.AddCommand(gitAutomationStatesListCmd)
	gitAutomationStatesCmd.AddCommand(gitAutomationStatesCreateCmd)
	gitAutomationStatesCmd.AddCommand(gitAutomationStatesUpdateCmd)
	gitAutomationStatesCmd.AddCommand(gitAutomationStatesDeleteCmd)

	gitAutomationBranchesCmd.AddCommand(gitAutomationBranchesCreateCmd)
	gitAutomationBranchesCmd.AddCommand(gitAutomationBranchesUpdateCmd)
	gitAutomationBranchesCmd.AddCommand(gitAutomationBranchesDeleteCmd)

	gitAutomationCmd.AddCommand(gitAutomationStatesCmd)
	gitAutomationCmd.AddCommand(gitAutomationBranchesCmd)

	rootCmd.AddCommand(gitAutomationCmd)
}
