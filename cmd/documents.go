package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var documentsCmd = &cobra.Command{
	Use:   "documents",
	Short: "Управление документами",
}

var documentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список документов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		docs, err := c.ListDocuments()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(docs)
		}

		if len(docs) == 0 {
			fmt.Fprintln(out, "Документы не найдены.")
			return nil
		}

		headers := []string{"ID", "TITLE", "PROJECT", "CREATOR"}
		rows := make([][]string, 0, len(docs))
		for _, d := range docs {
			projectName := ""
			if d.Project != nil {
				projectName = d.Project.Name
			}
			creatorName := ""
			if d.Creator != nil {
				creatorName = d.Creator.Name
			}
			rows = append(rows, []string{d.ID, d.Title, projectName, creatorName})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var documentsViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр документа по ID",
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
		doc, err := c.GetDocument(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(doc)
		}

		fmt.Fprintf(out, "ID:      %s\n", doc.ID)
		fmt.Fprintf(out, "Заголовок: %s\n", doc.Title)
		if doc.Project != nil {
			fmt.Fprintf(out, "Проект:  %s\n", doc.Project.Name)
		}
		if doc.Creator != nil {
			fmt.Fprintf(out, "Автор:   %s\n", doc.Creator.Name)
		}
		if doc.Content != "" {
			fmt.Fprintf(out, "\n%s\n", doc.Content)
		}
		return nil
	},
}

var documentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать документ",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		projectID, _ := cmd.Flags().GetString("project")

		if title == "" {
			return fmt.Errorf("необходимо указать --title")
		}

		c := newLinearClient(t)
		doc, err := c.CreateDocument(title, content, projectID)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(doc)
		}
		fmt.Fprintf(out, "Документ создан: %s (%s)\n", doc.Title, doc.ID)
		return nil
	},
}

var documentsUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить документ по ID",
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
		if title, _ := cmd.Flags().GetString("title"); title != "" {
			input["title"] = title
		}
		if content, _ := cmd.Flags().GetString("content"); content != "" {
			input["content"] = content
		}

		if len(input) == 0 {
			return fmt.Errorf("необходимо указать хотя бы один флаг для обновления")
		}

		c := newLinearClient(t)
		doc, err := c.UpdateDocument(args[0], input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(doc)
		}
		fmt.Fprintf(out, "Документ обновлён: %s (%s)\n", doc.Title, doc.ID)
		return nil
	},
}

var documentsDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить документ по ID",
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
		if err := c.DeleteDocument(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Документ удалён: %s\n", args[0])
		return nil
	},
}

var documentsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Поиск документов",
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
		docs, err := c.SearchDocuments(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(docs)
		}

		if len(docs) == 0 {
			fmt.Fprintln(out, "Документы не найдены.")
			return nil
		}

		headers := []string{"ID", "TITLE", "PROJECT", "CREATOR"}
		rows := make([][]string, 0, len(docs))
		for _, d := range docs {
			projectName := ""
			if d.Project != nil {
				projectName = d.Project.Name
			}
			creatorName := ""
			if d.Creator != nil {
				creatorName = d.Creator.Name
			}
			rows = append(rows, []string{d.ID, d.Title, projectName, creatorName})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

func init() {
	documentsCmd.AddCommand(documentsListCmd)
	documentsCmd.AddCommand(documentsViewCmd)

	documentsCreateCmd.Flags().String("title", "", "Заголовок документа (обязательно)")
	documentsCreateCmd.Flags().String("content", "", "Содержимое документа")
	documentsCreateCmd.Flags().String("project", "", "ID проекта")
	documentsCmd.AddCommand(documentsCreateCmd)

	documentsUpdateCmd.Flags().String("title", "", "Новый заголовок")
	documentsUpdateCmd.Flags().String("content", "", "Новое содержимое")
	documentsCmd.AddCommand(documentsUpdateCmd)

	documentsCmd.AddCommand(documentsDeleteCmd)
	documentsCmd.AddCommand(documentsSearchCmd)

	rootCmd.AddCommand(documentsCmd)
}
