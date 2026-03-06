package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueArchiveCmd = &cobra.Command{
	Use:   "archive <ISSUE-ID>",
	Short: "Архивировать задачу",
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
		if err := c.ArchiveIssue(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Задача архивирована: %s\n", args[0])
		return nil
	},
}

var issueUnarchiveCmd = &cobra.Command{
	Use:   "unarchive <ISSUE-ID>",
	Short: "Разархивировать задачу",
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
		if err := c.UnarchiveIssue(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Задача разархивирована: %s\n", args[0])
		return nil
	},
}

var issueDeleteCmd = &cobra.Command{
	Use:   "delete <ISSUE-ID>",
	Short: "Удалить задачу",
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
		if err := c.DeleteIssue(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Задача удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	issuesCmd.AddCommand(issueArchiveCmd)
	issuesCmd.AddCommand(issueUnarchiveCmd)
	issuesCmd.AddCommand(issueDeleteCmd)
}
