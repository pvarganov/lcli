package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var commentBody string

var issueCommentCmd = &cobra.Command{
	Use:   "comment <ID>",
	Short: "Добавить комментарий к задаче",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		comment, err := c.CreateComment(args[0], commentBody)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Комментарий добавлен (id: %s)\n", comment.ID)
		return nil
	},
}

var issueCommentsCmd = &cobra.Command{
	Use:   "comments <ID>",
	Short: "Список комментариев к задаче",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t := GetToken()
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		comments, err := c.ListComments(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if len(comments) == 0 {
			fmt.Fprintln(out, "Комментариев нет")
			return nil
		}

		for _, cm := range comments {
			author := "unknown"
			if cm.User != nil {
				author = cm.User.DisplayName
			}
			fmt.Fprintf(out, "%s  %s\n", cm.CreatedAt.Format("2006-01-02 15:04"), author)
			fmt.Fprintln(out, strings.Repeat("─", 60))
			fmt.Fprintln(out, cm.Body)
			fmt.Fprintln(out)
		}
		return nil
	},
}

func init() {
	issueCommentCmd.Flags().StringVar(&commentBody, "body", "", "Текст комментария (обязательно)")
	_ = issueCommentCmd.MarkFlagRequired("body")

	issueCmd.AddCommand(issueCommentCmd)
	issueCmd.AddCommand(issueCommentsCmd)
}
