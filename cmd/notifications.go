package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Управление уведомлениями",
}

var notificationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список уведомлений",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		limit, _ := cmd.Flags().GetInt("limit")
		after, _ := cmd.Flags().GetString("after")

		c := newLinearClient(t)
		notifications, pageInfo, err := c.ListNotifications(limit, after)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(notifications)
		}

		if len(notifications) == 0 {
			fmt.Fprintln(out, "Уведомления не найдены.")
			return nil
		}

		headers := []string{"ID", "TYPE", "READ", "CREATED", "ISSUE"}
		rows := make([][]string, 0, len(notifications))
		for _, n := range notifications {
			readAt := "no"
			if n.ReadAt != nil {
				readAt = "yes"
			}
			issueRef := ""
			if n.Issue != nil {
				issueRef = n.Issue.Identifier + " " + format.StripControlChars(n.Issue.Title)
				if len(issueRef) > 40 {
					issueRef = issueRef[:40] + "..."
				}
			}
			rows = append(rows, []string{
				n.ID,
				format.StripControlChars(n.Type),
				readAt,
				n.CreatedAt.Format("2006-01-02"),
				issueRef,
			})
		}
		format.TableWriter(out, headers, rows)

		if pageInfo != nil && pageInfo.HasNextPage {
			fmt.Fprintf(out, "\nЕсть ещё страницы. Используйте --after %s\n", pageInfo.EndCursor)
		}
		return nil
	},
}

var notificationsUnreadCountCmd = &cobra.Command{
	Use:   "unread-count",
	Short: "Количество непрочитанных уведомлений",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		count, err := c.GetNotificationsUnreadCount()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(map[string]int{"unreadCount": count})
		}

		fmt.Fprintf(out, "Непрочитанных уведомлений: %d\n", count)
		return nil
	},
}

var notificationsMarkReadCmd = &cobra.Command{
	Use:   "mark-read [ID|all]",
	Short: "Отметить уведомление(я) как прочитанные",
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
		out := cmd.OutOrStdout()

		if strings.ToLower(args[0]) == "all" {
			if err := c.MarkAllNotificationsRead(); err != nil {
				return err
			}
			fmt.Fprintln(out, "Все уведомления отмечены как прочитанные.")
			return nil
		}

		if err := c.MarkNotificationRead(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(out, "Уведомление %s отмечено как прочитанное.\n", args[0])
		return nil
	},
}

var notificationsArchiveCmd = &cobra.Command{
	Use:   "archive <ID>",
	Short: "Архивировать уведомление по ID",
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
		if err := c.ArchiveNotification(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Уведомление %s архивировано.\n", args[0])
		return nil
	},
}

func init() {
	notificationsListCmd.Flags().Int("limit", 50, "Максимальное количество уведомлений")
	notificationsListCmd.Flags().String("after", "", "Курсор для пагинации")

	notificationsCmd.AddCommand(notificationsListCmd)
	notificationsCmd.AddCommand(notificationsUnreadCountCmd)
	notificationsCmd.AddCommand(notificationsMarkReadCmd)
	notificationsCmd.AddCommand(notificationsArchiveCmd)
	rootCmd.AddCommand(notificationsCmd)
}
