package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var projectMilestonesCmd = &cobra.Command{
	Use:   "milestones",
	Short: "Управление вехами проекта",
}

var projectMilestonesListCmd = &cobra.Command{
	Use:   "list <PROJECT-ID>",
	Short: "Список вех проекта",
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
		milestones, err := c.ListProjectMilestones(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(milestones)
		}

		if len(milestones) == 0 {
			fmt.Fprintln(out, "Вехи не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "TARGET DATE", "DESCRIPTION"}
		rows := make([][]string, 0, len(milestones))
		for _, m := range milestones {
			desc := format.StripControlChars(m.Description)
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			rows = append(rows, []string{m.ID, format.StripControlChars(m.Name), m.TargetDate, desc})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var projectMilestonesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать веху проекта",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		projectID, _ := cmd.Flags().GetString("project-id")
		name, _ := cmd.Flags().GetString("name")
		targetDate, _ := cmd.Flags().GetString("target-date")
		description, _ := cmd.Flags().GetString("description")

		if projectID == "" {
			return fmt.Errorf("требуется --project-id")
		}
		if name == "" {
			return fmt.Errorf("требуется --name")
		}

		c := newLinearClient(t)
		ms, err := c.CreateProjectMilestone(client.CreateProjectMilestoneInput{
			ProjectID:   projectID,
			Name:        name,
			TargetDate:  targetDate,
			Description: description,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Веха создана: %s (%s)\n", ms.Name, ms.ID)
		return nil
	},
}

var projectMilestonesUpdateCmd = &cobra.Command{
	Use:   "update <MILESTONE-ID>",
	Short: "Обновить веху проекта",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		targetDate, _ := cmd.Flags().GetString("target-date")
		description, _ := cmd.Flags().GetString("description")

		c := newLinearClient(t)
		ms, err := c.UpdateProjectMilestone(args[0], client.UpdateProjectMilestoneInput{
			Name:        name,
			TargetDate:  targetDate,
			Description: description,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Веха обновлена: %s (%s)\n", ms.Name, ms.ID)
		return nil
	},
}

var projectMilestonesDeleteCmd = &cobra.Command{
	Use:   "delete <MILESTONE-ID>",
	Short: "Удалить веху проекта",
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
		if err := c.DeleteProjectMilestone(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Веха удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	projectMilestonesCreateCmd.Flags().String("project-id", "", "ID проекта (обязательно)")
	projectMilestonesCreateCmd.Flags().String("name", "", "Название вехи (обязательно)")
	projectMilestonesCreateCmd.Flags().String("target-date", "", "Целевая дата (YYYY-MM-DD)")
	projectMilestonesCreateCmd.Flags().String("description", "", "Описание вехи")

	projectMilestonesUpdateCmd.Flags().String("name", "", "Новое название вехи")
	projectMilestonesUpdateCmd.Flags().String("target-date", "", "Новая целевая дата (YYYY-MM-DD)")
	projectMilestonesUpdateCmd.Flags().String("description", "", "Новое описание")

	projectMilestonesCmd.AddCommand(projectMilestonesListCmd)
	projectMilestonesCmd.AddCommand(projectMilestonesCreateCmd)
	projectMilestonesCmd.AddCommand(projectMilestonesUpdateCmd)
	projectMilestonesCmd.AddCommand(projectMilestonesDeleteCmd)
	projectsCmd.AddCommand(projectMilestonesCmd)
}
