package cmd

import (
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/spf13/cobra"
)

var projectRelationsCmd = &cobra.Command{
	Use:   "relations",
	Short: "Управление связями проектов",
}

var projectRelationsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавить связь между проектами",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		projectID, _ := cmd.Flags().GetString("project")
		relatedID, _ := cmd.Flags().GetString("related")
		relType, _ := cmd.Flags().GetString("type")

		c := newLinearClient(t)
		rel, err := c.CreateProjectRelation(client.CreateProjectRelationInput{
			ProjectID:        projectID,
			RelatedProjectID: relatedID,
			Type:             relType,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Связь создана: %s (%s → %s)\n", rel.ID, rel.Project.Name, rel.RelatedProject.Name)
		return nil
	},
}

var projectRelationsRemoveCmd = &cobra.Command{
	Use:   "remove <RELATION-ID>",
	Short: "Удалить связь между проектами",
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
		if err := c.DeleteProjectRelation(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Связь удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	projectRelationsAddCmd.Flags().String("project", "", "ID проекта (обязательно)")
	projectRelationsAddCmd.Flags().String("related", "", "ID связанного проекта (обязательно)")
	projectRelationsAddCmd.Flags().String("type", "related", "Тип связи: related, blocks, blocked_by")
	_ = projectRelationsAddCmd.MarkFlagRequired("project")
	_ = projectRelationsAddCmd.MarkFlagRequired("related")

	projectRelationsCmd.AddCommand(projectRelationsAddCmd)
	projectRelationsCmd.AddCommand(projectRelationsRemoveCmd)

	projectsCmd.AddCommand(projectRelationsCmd)
}
