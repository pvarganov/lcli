package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var commentBody string
var commentParentID string

var issueCommentCmd = &cobra.Command{
	Use:   "comment <ID>",
	Short: "Добавить комментарий к задаче",
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
		comment, err := c.CreateComment(client.CreateCommentInput{
			IssueID:  args[0],
			Body:     commentBody,
			ParentID: commentParentID,
		})
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
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		comments, err := c.ListComments(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(comments)
		}

		if len(comments) == 0 {
			fmt.Fprintln(out, "Комментариев нет")
			return nil
		}

		for _, cm := range comments {
			author := "unknown"
			if cm.User != nil {
				author = format.StripControlChars(cm.User.DisplayName)
			}
			fmt.Fprintf(out, "%s  %s\n", cm.CreatedAt.Format("2006-01-02 15:04"), author)
			fmt.Fprintln(out, strings.Repeat("─", 60))
			fmt.Fprintln(out, format.StripControlChars(cm.Body))
			fmt.Fprintln(out)
		}
		return nil
	},
}

var issueCommentUpdateCmd = &cobra.Command{
	Use:   "comment-update <COMMENT-ID>",
	Short: "Обновить комментарий по ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		body, _ := cmd.Flags().GetString("body")
		c := newLinearClient(t)
		comment, err := c.UpdateComment(args[0], body)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Комментарий обновлён (id: %s)\n", comment.ID)
		return nil
	},
}

var issueCommentDeleteCmd = &cobra.Command{
	Use:   "comment-delete <COMMENT-ID>",
	Short: "Удалить комментарий по ID",
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
		if err := c.DeleteComment(args[0]); err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Комментарий удалён (id: %s)\n", args[0])
		return nil
	},
}

var issueCommentResolveCmd = &cobra.Command{
	Use:   "comment-resolve <COMMENT-ID>",
	Short: "Пометить комментарий как разрешённый",
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
		comment, err := c.ResolveComment(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Комментарий помечен как разрешённый (id: %s)\n", comment.ID)
		return nil
	},
}

var issueCommentUnresolveCmd = &cobra.Command{
	Use:   "comment-unresolve <COMMENT-ID>",
	Short: "Снять пометку разрешения с комментария",
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
		comment, err := c.UnresolveComment(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Пометка разрешения снята (id: %s)\n", comment.ID)
		return nil
	},
}

func init() {
	issueCommentCmd.Flags().StringVar(&commentBody, "body", "", "Текст комментария (обязательно)")
	_ = issueCommentCmd.MarkFlagRequired("body")
	issueCommentCmd.Flags().StringVar(&commentParentID, "parent-id", "", "ID родительского комментария (для вложенного ответа)")

	issueCommentUpdateCmd.Flags().String("body", "", "Новый текст комментария (обязательно)")
	_ = issueCommentUpdateCmd.MarkFlagRequired("body")

	issueCmd.AddCommand(issueCommentCmd)
	issueCmd.AddCommand(issueCommentsCmd)
	issueCmd.AddCommand(issueCommentUpdateCmd)
	issueCmd.AddCommand(issueCommentDeleteCmd)
	issueCmd.AddCommand(issueCommentResolveCmd)
	issueCmd.AddCommand(issueCommentUnresolveCmd)
}
