package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var favoritesCmd = &cobra.Command{
	Use:   "favorites",
	Short: "Управление избранными элементами",
}

var favoritesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список избранных элементов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		favs, err := c.ListFavorites()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(favs)
		}

		if len(favs) == 0 {
			fmt.Fprintln(out, "Избранные элементы не найдены.")
			return nil
		}

		headers := []string{"ID", "TYPE", "ENTITY"}
		rows := make([][]string, 0, len(favs))
		for _, f := range favs {
			entity := ""
			switch f.Type {
			case "issue":
				if f.Issue != nil {
					entity = fmt.Sprintf("%s: %s", f.Issue.Identifier, f.Issue.Title)
				}
			case "project":
				if f.Project != nil {
					entity = f.Project.Name
				}
			case "label":
				if f.Label != nil {
					entity = f.Label.Name
				}
			case "customView":
				if f.CustomView != nil {
					entity = f.CustomView.Name
				}
			}
			rows = append(rows, []string{f.ID, f.Type, entity})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var favoritesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавить элемент в избранное",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		entityType, _ := cmd.Flags().GetString("type")
		entityID, _ := cmd.Flags().GetString("id")

		if entityType == "" {
			return fmt.Errorf("необходимо указать --type (issue, project, label, customView)")
		}
		if entityID == "" {
			return fmt.Errorf("необходимо указать --id")
		}

		c := newLinearClient(t)
		fav, err := c.CreateFavorite(entityType, entityID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(fav)
		}
		fmt.Fprintf(out, "Добавлено в избранное: %s (%s)\n", fav.Type, fav.ID)
		return nil
	},
}

var favoritesRemoveCmd = &cobra.Command{
	Use:   "remove <ID>",
	Short: "Удалить элемент из избранного по ID",
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
		if err := c.DeleteFavorite(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Удалено из избранного: %s\n", args[0])
		return nil
	},
}

func init() {
	favoritesAddCmd.Flags().String("type", "", "Тип элемента: issue, project, label, customView (обязательно)")
	favoritesAddCmd.Flags().String("id", "", "ID элемента (обязательно)")

	favoritesCmd.AddCommand(favoritesListCmd)
	favoritesCmd.AddCommand(favoritesAddCmd)
	favoritesCmd.AddCommand(favoritesRemoveCmd)

	rootCmd.AddCommand(favoritesCmd)
}
