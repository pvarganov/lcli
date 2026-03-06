package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var roadmapsCmd = &cobra.Command{
	Use:   "roadmaps",
	Short: "Управление дорожными картами",
}

var roadmapsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список дорожных карт",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		roadmaps, err := c.ListRoadmaps()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(roadmaps)
		}

		if len(roadmaps) == 0 {
			fmt.Fprintln(out, "Дорожные карты не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "OWNER"}
		rows := make([][]string, 0, len(roadmaps))
		for _, r := range roadmaps {
			ownerName := ""
			if r.Owner != nil {
				ownerName = r.Owner.Name
			}
			rows = append(rows, []string{r.ID, r.Name, ownerName})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var roadmapsViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр дорожной карты по ID",
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
		rm, err := c.GetRoadmap(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(rm)
		}

		fmt.Fprintf(out, "ID:      %s\n", rm.ID)
		fmt.Fprintf(out, "Название: %s\n", rm.Name)
		if rm.Owner != nil {
			fmt.Fprintf(out, "Владелец: %s\n", rm.Owner.Name)
		}
		if rm.Description != "" {
			fmt.Fprintf(out, "\n%s\n", rm.Description)
		}
		return nil
	},
}

var roadmapsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать дорожную карту",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		ownerID, _ := cmd.Flags().GetString("owner")

		if name == "" {
			return fmt.Errorf("необходимо указать --name")
		}

		c := newLinearClient(t)
		rm, err := c.CreateRoadmap(name, description, ownerID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(rm)
		}
		fmt.Fprintf(out, "Дорожная карта создана: %s (%s)\n", rm.Name, rm.ID)
		return nil
	},
}

var roadmapsUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить дорожную карту по ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		input := map[string]any{}
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			input["name"] = name
		}
		if description, _ := cmd.Flags().GetString("description"); description != "" {
			input["description"] = description
		}

		if len(input) == 0 {
			return fmt.Errorf("необходимо указать хотя бы один флаг для обновления")
		}

		c := newLinearClient(t)
		rm, err := c.UpdateRoadmap(args[0], input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(rm)
		}
		fmt.Fprintf(out, "Дорожная карта обновлена: %s (%s)\n", rm.Name, rm.ID)
		return nil
	},
}

var roadmapsDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить дорожную карту по ID",
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
		if err := c.DeleteRoadmap(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Дорожная карта удалена: %s\n", args[0])
		return nil
	},
}

var roadmapsAddProjectCmd = &cobra.Command{
	Use:   "add-project",
	Short: "Добавить проект в дорожную карту",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		roadmapID, _ := cmd.Flags().GetString("roadmap")
		projectID, _ := cmd.Flags().GetString("project")

		if roadmapID == "" {
			return fmt.Errorf("требуется --roadmap")
		}
		if projectID == "" {
			return fmt.Errorf("требуется --project")
		}

		c := newLinearClient(t)
		rel, err := c.CreateRoadmapToProject(roadmapID, projectID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(rel)
		}
		fmt.Fprintf(out, "Проект добавлен в дорожную карту: %s (%s → %s)\n", rel.ID, rel.Roadmap.Name, rel.Project.Name)
		return nil
	},
}

var roadmapsRemoveProjectCmd = &cobra.Command{
	Use:   "remove-project <RELATION-ID>",
	Short: "Удалить проект из дорожной карты",
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
		if err := c.DeleteRoadmapToProject(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Проект удалён из дорожной карты: %s\n", args[0])
		return nil
	},
}

func init() {
	roadmapsCreateCmd.Flags().String("name", "", "Название дорожной карты (обязательно)")
	roadmapsCreateCmd.Flags().String("description", "", "Описание")
	roadmapsCreateCmd.Flags().String("owner", "", "ID владельца")

	roadmapsUpdateCmd.Flags().String("name", "", "Новое название")
	roadmapsUpdateCmd.Flags().String("description", "", "Новое описание")

	roadmapsAddProjectCmd.Flags().String("roadmap", "", "ID дорожной карты (обязательно)")
	roadmapsAddProjectCmd.Flags().String("project", "", "ID проекта (обязательно)")

	roadmapsCmd.AddCommand(roadmapsListCmd)
	roadmapsCmd.AddCommand(roadmapsViewCmd)
	roadmapsCmd.AddCommand(roadmapsCreateCmd)
	roadmapsCmd.AddCommand(roadmapsUpdateCmd)
	roadmapsCmd.AddCommand(roadmapsDeleteCmd)
	roadmapsCmd.AddCommand(roadmapsAddProjectCmd)
	roadmapsCmd.AddCommand(roadmapsRemoveProjectCmd)
	rootCmd.AddCommand(roadmapsCmd)
}
