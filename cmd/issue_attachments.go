package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var issueAttachmentsCmd = &cobra.Command{
	Use:   "attachments",
	Short: "Управление вложениями задачи",
}

var issueAttachmentsListCmd = &cobra.Command{
	Use:   "list <ISSUE-ID>",
	Short: "Список вложений задачи",
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
		attachments, err := c.ListAttachments(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(attachments)
		}

		if len(attachments) == 0 {
			fmt.Fprintln(out, "Вложения не найдены.")
			return nil
		}

		headers := []string{"ID", "TITLE", "URL", "TYPE", "SUBTITLE"}
		rows := make([][]string, 0, len(attachments))
		for _, att := range attachments {
			rows = append(rows, []string{
				att.ID,
				att.Title,
				att.URL,
				att.SourceType,
				att.Subtitle,
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var issueAttachmentsLinkURLCmd = &cobra.Command{
	Use:   "link-url <ISSUE-ID>",
	Short: "Привязать URL к задаче как вложение",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		url, _ := cmd.Flags().GetString("url")
		title, _ := cmd.Flags().GetString("title")

		if url == "" {
			return fmt.Errorf("необходимо указать --url")
		}

		c := newLinearClient(t)
		att, err := c.AttachmentLinkURL(args[0], url, title)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(att)
		}
		fmt.Fprintf(out, "Вложение создано: %s (%s)\n", att.Title, att.ID)
		return nil
	},
}

var issueAttachmentsLinkGitHubPRCmd = &cobra.Command{
	Use:   "link-github-pr <ISSUE-ID>",
	Short: "Привязать GitHub PR к задаче",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		url, _ := cmd.Flags().GetString("url")
		title, _ := cmd.Flags().GetString("title")

		if url == "" {
			return fmt.Errorf("необходимо указать --url")
		}

		c := newLinearClient(t)
		att, err := c.AttachmentLinkGitHubPR(args[0], url, title)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(att)
		}
		fmt.Fprintf(out, "GitHub PR привязан: %s (%s)\n", att.Title, att.ID)
		return nil
	},
}

var issueAttachmentsDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить вложение по ID",
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
		if err := c.AttachmentDelete(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Вложение удалено: %s\n", args[0])
		return nil
	},
}

func init() {
	issueAttachmentsCmd.AddCommand(issueAttachmentsListCmd)

	issueAttachmentsLinkURLCmd.Flags().String("url", "", "URL для привязки (обязательно)")
	issueAttachmentsLinkURLCmd.Flags().String("title", "", "Заголовок вложения")
	issueAttachmentsCmd.AddCommand(issueAttachmentsLinkURLCmd)

	issueAttachmentsLinkGitHubPRCmd.Flags().String("url", "", "URL GitHub PR (обязательно)")
	issueAttachmentsLinkGitHubPRCmd.Flags().String("title", "", "Заголовок вложения")
	issueAttachmentsCmd.AddCommand(issueAttachmentsLinkGitHubPRCmd)

	issueAttachmentsCmd.AddCommand(issueAttachmentsDeleteCmd)

	issuesCmd.AddCommand(issueAttachmentsCmd)
}

