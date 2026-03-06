package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueSubscribeCmd = &cobra.Command{
	Use:   "subscribe <ISSUE-ID>",
	Short: "Подписаться на задачу",
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
		if err := c.SubscribeToIssue(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Подписка оформлена: %s\n", args[0])
		return nil
	},
}

var issueUnsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe <ISSUE-ID>",
	Short: "Отписаться от задачи",
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
		if err := c.UnsubscribeFromIssue(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Подписка отменена: %s\n", args[0])
		return nil
	},
}

func init() {
	issuesCmd.AddCommand(issueSubscribeCmd)
	issuesCmd.AddCommand(issueUnsubscribeCmd)
}
