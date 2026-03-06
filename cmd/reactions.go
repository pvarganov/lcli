package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commentReactCmd = &cobra.Command{
	Use:   "react <COMMENT-ID> <emoji>",
	Short: "Добавить реакцию на комментарий",
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
		id, err := c.CreateReaction(args[0], args[1])
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Реакция добавлена (id: %s)\n", id)
		return nil
	},
}

var commentUnreactCmd = &cobra.Command{
	Use:   "unreact <REACTION-ID>",
	Short: "Удалить реакцию по ID",
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
		if err := c.DeleteReaction(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Реакция удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	issuesCmd.AddCommand(commentReactCmd)
	issuesCmd.AddCommand(commentUnreactCmd)
}
