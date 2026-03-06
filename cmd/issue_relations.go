package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var relationsCmd = &cobra.Command{
	Use:   "relations",
	Short: "Управление связями задач",
}

var relationsListCmd = &cobra.Command{
	Use:   "list <ISSUE-ID>",
	Short: "Список связей задачи",
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
		rels, err := c.ListIssueRelations(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(rels)
		}

		if len(rels) == 0 {
			fmt.Fprintln(out, "Связи не найдены.")
			return nil
		}

		headers := []string{"ID", "TYPE", "RELATED ISSUE", "TITLE"}
		rows := make([][]string, 0, len(rels))
		for _, r := range rels {
			rows = append(rows, []string{
				r.ID,
				r.Type,
				r.RelatedIssue.Identifier,
				format.StripControlChars(r.RelatedIssue.Title),
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var (
	relationIssueID        string
	relationRelatedIssueID string
	relationRelType        string
)

var relationsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавить связь между задачами",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		rel, err := c.CreateIssueRelation(client.CreateIssueRelationInput{
			IssueID:        relationIssueID,
			RelatedIssueID: relationRelatedIssueID,
			Type:           relationRelType,
		})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Связь создана: %s (%s → %s)\n", rel.ID, relationRelType, rel.RelatedIssue.Identifier)
		return nil
	},
}

var relationsRemoveCmd = &cobra.Command{
	Use:   "remove <RELATION-ID>",
	Short: "Удалить связь между задачами",
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
		if err := c.DeleteIssueRelation(args[0]); err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Связь удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	relationsAddCmd.Flags().StringVar(&relationIssueID, "issue", "", "ID задачи (обязательно)")
	relationsAddCmd.Flags().StringVar(&relationRelatedIssueID, "related", "", "ID связанной задачи (обязательно)")
	relationsAddCmd.Flags().StringVar(&relationRelType, "type", "related", "Тип связи: blocks, blocked_by, related, duplicate, duplicate_of")
	_ = relationsAddCmd.MarkFlagRequired("issue")
	_ = relationsAddCmd.MarkFlagRequired("related")

	relationsCmd.AddCommand(relationsListCmd)
	relationsCmd.AddCommand(relationsAddCmd)
	relationsCmd.AddCommand(relationsRemoveCmd)

	issuesCmd.AddCommand(relationsCmd)
}
