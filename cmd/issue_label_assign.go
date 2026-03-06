package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Добавление и удаление меток на задачах",
}

var labelAddCmd = &cobra.Command{
	Use:   "add <ISSUE-ID> <label-name>",
	Short: "Добавить метку к задаче (по имени метки)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		issueID := args[0]
		labelName := args[1]

		c := newLinearClient(t)

		labels, err := c.ListIssueLabels()
		if err != nil {
			return err
		}

		var labelID string
		for _, l := range labels {
			if l.Name == labelName {
				labelID = l.ID
				break
			}
		}
		if labelID == "" {
			return fmt.Errorf("метка не найдена: %s", labelName)
		}

		if err := c.AddLabelToIssue(issueID, labelID); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Метка \"%s\" добавлена к задаче %s\n", labelName, issueID)
		return nil
	},
}

var labelRemoveCmd = &cobra.Command{
	Use:   "remove <ISSUE-ID> <label-name>",
	Short: "Удалить метку из задачи (по имени метки)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		issueID := args[0]
		labelName := args[1]

		c := newLinearClient(t)

		labels, err := c.ListIssueLabels()
		if err != nil {
			return err
		}

		var labelID string
		for _, l := range labels {
			if l.Name == labelName {
				labelID = l.ID
				break
			}
		}
		if labelID == "" {
			return fmt.Errorf("метка не найдена: %s", labelName)
		}

		if err := c.RemoveLabelFromIssue(issueID, labelID); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Метка \"%s\" удалена из задачи %s\n", labelName, issueID)
		return nil
	},
}

func init() {
	labelCmd.AddCommand(labelAddCmd)
	labelCmd.AddCommand(labelRemoveCmd)
	issuesCmd.AddCommand(labelCmd)
}
