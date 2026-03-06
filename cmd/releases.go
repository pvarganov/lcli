package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var releasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "Управление релизами [ALPHA]",
}

var releasesPipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Управление пайплайнами релизов",
}

var releasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список релизов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		releases, err := c.ListReleases()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(releases)
		}

		if len(releases) == 0 {
			fmt.Fprintln(out, "Релизы не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "PIPELINE", "STAGE", "COMPLETED_AT"}
		rows := make([][]string, 0, len(releases))
		for _, r := range releases {
			pipeline := ""
			if r.Pipeline != nil {
				pipeline = r.Pipeline.Name
			}
			stage := ""
			if r.Stage != nil {
				stage = r.Stage.Name
			}
			rows = append(rows, []string{r.ID, r.Name, pipeline, stage, r.CompletedAt})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var releasesViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр релиза",
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
		release, err := c.GetRelease(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(release)
		}

		fmt.Fprintf(out, "ID:          %s\n", release.ID)
		fmt.Fprintf(out, "Name:        %s\n", release.Name)
		fmt.Fprintf(out, "Description: %s\n", release.Description)
		fmt.Fprintf(out, "CommitSHA:   %s\n", release.CommitSha)
		fmt.Fprintf(out, "Created:     %s\n", release.CreatedAt)
		fmt.Fprintf(out, "Completed:   %s\n", release.CompletedAt)
		if release.Pipeline != nil {
			fmt.Fprintf(out, "Pipeline:    %s (%s)\n", release.Pipeline.Name, release.Pipeline.ID)
		}
		if release.Stage != nil {
			fmt.Fprintf(out, "Stage:       %s\n", release.Stage.Name)
		}
		return nil
	},
}

var releasesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать релиз",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("флаг --name обязателен")
		}
		pipelineID, _ := cmd.Flags().GetString("pipeline")
		description, _ := cmd.Flags().GetString("description")
		version, _ := cmd.Flags().GetString("version")

		c := newLinearClient(t)
		release, err := c.CreateRelease(name, pipelineID, description, version)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(release)
		}

		fmt.Fprintf(out, "Релиз создан: %s (%s)\n", release.Name, release.ID)
		return nil
	},
}

var releasesCompleteCmd = &cobra.Command{
	Use:   "complete <PIPELINE-ID>",
	Short: "Завершить релиз для пайплайна",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		version, _ := cmd.Flags().GetString("version")
		commitSha, _ := cmd.Flags().GetString("commit-sha")

		c := newLinearClient(t)
		release, err := c.CompleteRelease(args[0], version, commitSha)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(release)
		}

		fmt.Fprintf(out, "Релиз завершён: %s (%s)\n", release.Name, release.ID)
		return nil
	},
}

var releasesDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить релиз",
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
		if err := c.DeleteRelease(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Релиз удалён: %s\n", args[0])
		return nil
	},
}

var releasesSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Поиск релизов",
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
		releases, err := c.SearchReleases(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(releases)
		}

		if len(releases) == 0 {
			fmt.Fprintln(out, "Релизы не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "PIPELINE", "STAGE"}
		rows := make([][]string, 0, len(releases))
		for _, r := range releases {
			pipeline := ""
			if r.Pipeline != nil {
				pipeline = r.Pipeline.Name
			}
			stage := ""
			if r.Stage != nil {
				stage = r.Stage.Name
			}
			rows = append(rows, []string{r.ID, r.Name, pipeline, stage})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var releasesPipelinesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список пайплайнов релизов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		pipelines, err := c.ListReleasePipelines()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(pipelines)
		}

		if len(pipelines) == 0 {
			fmt.Fprintln(out, "Пайплайны не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "SLUG", "TYPE"}
		rows := make([][]string, 0, len(pipelines))
		for _, p := range pipelines {
			rows = append(rows, []string{p.ID, p.Name, p.SlugID, p.Type})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var releasesPipelinesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать пайплайн релизов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("флаг --name обязателен")
		}

		c := newLinearClient(t)
		pipeline, err := c.CreateReleasePipeline(name)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(pipeline)
		}

		fmt.Fprintf(out, "Пайплайн создан: %s (%s)\n", pipeline.Name, pipeline.ID)
		return nil
	},
}

var releasesPipelinesDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить пайплайн релизов",
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
		if err := c.DeleteReleasePipeline(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Пайплайн удалён: %s\n", args[0])
		return nil
	},
}

func init() {
	releasesCreateCmd.Flags().String("name", "", "Название релиза (обязательно)")
	releasesCreateCmd.Flags().String("pipeline", "", "ID пайплайна")
	releasesCreateCmd.Flags().String("description", "", "Описание релиза")
	releasesCreateCmd.Flags().String("version", "", "Версия релиза")

	releasesCompleteCmd.Flags().String("version", "", "Версия релиза")
	releasesCompleteCmd.Flags().String("commit-sha", "", "SHA коммита")

	releasesPipelinesCreateCmd.Flags().String("name", "", "Название пайплайна (обязательно)")

	releasesCmd.AddCommand(releasesListCmd)
	releasesCmd.AddCommand(releasesViewCmd)
	releasesCmd.AddCommand(releasesCreateCmd)
	releasesCmd.AddCommand(releasesCompleteCmd)
	releasesCmd.AddCommand(releasesDeleteCmd)
	releasesCmd.AddCommand(releasesSearchCmd)
	releasesCmd.AddCommand(releasesPipelinesCmd)

	releasesPipelinesCmd.AddCommand(releasesPipelinesListCmd)
	releasesPipelinesCmd.AddCommand(releasesPipelinesCreateCmd)
	releasesPipelinesCmd.AddCommand(releasesPipelinesDeleteCmd)

	rootCmd.AddCommand(releasesCmd)
}
