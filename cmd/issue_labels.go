package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/client"
	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Управление метками задач",
}

var labelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список меток",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		labels, err := c.ListIssueLabels()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(labels)
		}

		if len(labels) == 0 {
			fmt.Fprintln(out, "Метки не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "COLOR", "DESCRIPTION"}
		rows := make([][]string, 0, len(labels))
		for _, l := range labels {
			rows = append(rows, []string{
				l.ID,
				format.StripControlChars(l.Name),
				l.Color,
				format.StripControlChars(l.Description),
			})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var (
	labelName        string
	labelColor       string
	labelDescription string
	labelTeam        string
)

var labelsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать метку",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)

		var teamID string
		if labelTeam != "" {
			team, err := c.GetTeamByKey(labelTeam)
			if err != nil {
				return err
			}
			if team == nil {
				return fmt.Errorf("команда не найдена: %s", labelTeam)
			}
			teamID = team.ID
		}

		label, err := c.CreateIssueLabel(client.CreateIssueLabelInput{
			Name:        labelName,
			Color:       labelColor,
			Description: labelDescription,
			TeamID:      teamID,
		})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Метка создана: %s (id: %s)\n", label.Name, label.ID)
		return nil
	},
}

var labelsUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить метку",
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
		label, err := c.UpdateIssueLabel(args[0], client.UpdateIssueLabelInput{
			Name:        labelName,
			Color:       labelColor,
			Description: labelDescription,
		})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Метка обновлена: %s (id: %s)\n", label.Name, label.ID)
		return nil
	},
}

var labelsDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить метку",
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
		if err := c.DeleteIssueLabel(args[0]); err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Метка удалена: %s\n", args[0])
		return nil
	},
}

func init() {
	labelsCreateCmd.Flags().StringVar(&labelName, "name", "", "Название метки (обязательно)")
	labelsCreateCmd.Flags().StringVar(&labelColor, "color", "#000000", "Цвет метки (hex)")
	labelsCreateCmd.Flags().StringVar(&labelDescription, "description", "", "Описание метки")
	labelsCreateCmd.Flags().StringVar(&labelTeam, "team", "", "Ключ команды (опционально)")
	_ = labelsCreateCmd.MarkFlagRequired("name")

	labelsUpdateCmd.Flags().StringVar(&labelName, "name", "", "Новое название метки")
	labelsUpdateCmd.Flags().StringVar(&labelColor, "color", "", "Новый цвет метки (hex)")
	labelsUpdateCmd.Flags().StringVar(&labelDescription, "description", "", "Новое описание метки")

	labelsCmd.AddCommand(labelsListCmd)
	labelsCmd.AddCommand(labelsCreateCmd)
	labelsCmd.AddCommand(labelsUpdateCmd)
	labelsCmd.AddCommand(labelsDeleteCmd)

	issuesCmd.AddCommand(labelsCmd)
}
