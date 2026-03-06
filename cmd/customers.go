package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pavelvarganov/lcli/internal/format"
	"github.com/spf13/cobra"
)

var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Управление клиентами (CRM)",
}

var customersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список клиентов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		customers, err := c.ListCustomers()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(customers)
		}

		if len(customers) == 0 {
			fmt.Fprintln(out, "Клиенты не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "STATUS", "TIER"}
		rows := make([][]string, 0, len(customers))
		for _, cu := range customers {
			statusName := ""
			if cu.Status != nil {
				statusName = cu.Status.DisplayName
			}
			tierName := ""
			if cu.Tier != nil {
				tierName = cu.Tier.DisplayName
			}
			rows = append(rows, []string{cu.ID, format.StripControlChars(cu.Name), format.StripControlChars(statusName), format.StripControlChars(tierName)})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var customersViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "Просмотр клиента по ID",
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
		cu, err := c.GetCustomer(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(cu)
		}

		fmt.Fprintf(out, "ID:       %s\n", cu.ID)
		fmt.Fprintf(out, "Название: %s\n", format.StripControlChars(cu.Name))
		if cu.Owner != nil {
			fmt.Fprintf(out, "Владелец: %s\n", format.StripControlChars(cu.Owner.Name))
		}
		if cu.Status != nil {
			fmt.Fprintf(out, "Статус:   %s\n", format.StripControlChars(cu.Status.DisplayName))
		}
		if cu.Tier != nil {
			fmt.Fprintf(out, "Уровень:  %s\n", format.StripControlChars(cu.Tier.DisplayName))
		}
		if cu.Revenue > 0 {
			fmt.Fprintf(out, "Доход:    %d\n", cu.Revenue)
		}
		return nil
	},
}

var customersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать клиента",
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
			return fmt.Errorf("необходимо указать --name")
		}

		input := map[string]any{}
		if logoURL, _ := cmd.Flags().GetString("logo-url"); logoURL != "" {
			input["logoUrl"] = logoURL
		}
		if ownerID, _ := cmd.Flags().GetString("owner"); ownerID != "" {
			input["ownerId"] = ownerID
		}
		if statusID, _ := cmd.Flags().GetString("status"); statusID != "" {
			input["statusId"] = statusID
		}
		if tierID, _ := cmd.Flags().GetString("tier"); tierID != "" {
			input["tierId"] = tierID
		}

		c := newLinearClient(t)
		cu, err := c.CreateCustomer(name, input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(cu)
		}
		fmt.Fprintf(out, "Клиент создан: %s (%s)\n", cu.Name, cu.ID)
		return nil
	},
}

var customersUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Обновить клиента по ID",
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
		if logoURL, _ := cmd.Flags().GetString("logo-url"); logoURL != "" {
			input["logoUrl"] = logoURL
		}
		if ownerID, _ := cmd.Flags().GetString("owner"); ownerID != "" {
			input["ownerId"] = ownerID
		}
		if statusID, _ := cmd.Flags().GetString("status"); statusID != "" {
			input["statusId"] = statusID
		}
		if tierID, _ := cmd.Flags().GetString("tier"); tierID != "" {
			input["tierId"] = tierID
		}

		if len(input) == 0 {
			return fmt.Errorf("необходимо указать хотя бы один флаг для обновления")
		}

		c := newLinearClient(t)
		cu, err := c.UpdateCustomer(args[0], input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(cu)
		}
		fmt.Fprintf(out, "Клиент обновлён: %s (%s)\n", cu.Name, cu.ID)
		return nil
	},
}

var customersDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Удалить клиента по ID",
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
		if err := c.DeleteCustomer(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Клиент удалён: %s\n", args[0])
		return nil
	},
}

var customersNeedsCmd = &cobra.Command{
	Use:   "needs",
	Short: "Управление потребностями клиентов",
}

var customersNeedsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список потребностей клиентов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		needs, err := c.ListCustomerNeeds()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(needs)
		}

		if len(needs) == 0 {
			fmt.Fprintln(out, "Потребности клиентов не найдены.")
			return nil
		}

		headers := []string{"ID", "BODY", "CUSTOMER", "PRIORITY"}
		rows := make([][]string, 0, len(needs))
		for _, n := range needs {
			customerName := ""
			if n.Customer != nil {
				customerName = n.Customer.Name
			}
			body := n.Body
			if len(body) > 50 {
				body = body[:47] + "..."
			}
			rows = append(rows, []string{n.ID, body, customerName, fmt.Sprintf("%.0f", n.Priority)})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var customersNeedsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать потребность клиента",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		body, _ := cmd.Flags().GetString("body")
		if body == "" {
			return fmt.Errorf("необходимо указать --body")
		}

		input := map[string]any{"body": body}
		if customerID, _ := cmd.Flags().GetString("customer"); customerID != "" {
			input["customerId"] = customerID
		}
		if issueID, _ := cmd.Flags().GetString("issue"); issueID != "" {
			input["issueId"] = issueID
		}

		c := newLinearClient(t)
		need, err := c.CreateCustomerNeed(input)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(need)
		}
		fmt.Fprintf(out, "Потребность создана: %s\n", need.ID)
		return nil
	},
}

var customersStatusesCmd = &cobra.Command{
	Use:   "statuses",
	Short: "Управление статусами клиентов",
}

var customersStatusesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список статусов клиентов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		statuses, err := c.ListCustomerStatuses()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(statuses)
		}

		if len(statuses) == 0 {
			fmt.Fprintln(out, "Статусы клиентов не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "DISPLAY NAME", "COLOR"}
		rows := make([][]string, 0, len(statuses))
		for _, s := range statuses {
			rows = append(rows, []string{s.ID, s.Name, s.DisplayName, s.Color})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

var customersTiersCmd = &cobra.Command{
	Use:   "tiers",
	Short: "Управление уровнями клиентов",
}

var customersTiersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список уровней клиентов",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := GetToken()
		if err != nil {
			return fmt.Errorf("ошибка загрузки токена: %w", err)
		}
		if t == "" {
			return fmt.Errorf("токен не настроен. Используйте --token или запустите `lcli auth login`")
		}

		c := newLinearClient(t)
		tiers, err := c.ListCustomerTiers()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()

		if GetOutputFormat() == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(tiers)
		}

		if len(tiers) == 0 {
			fmt.Fprintln(out, "Уровни клиентов не найдены.")
			return nil
		}

		headers := []string{"ID", "NAME", "DISPLAY NAME", "COLOR"}
		rows := make([][]string, 0, len(tiers))
		for _, tier := range tiers {
			rows = append(rows, []string{tier.ID, tier.Name, tier.DisplayName, tier.Color})
		}
		format.TableWriter(out, headers, rows)
		return nil
	},
}

func init() {
	customersCreateCmd.Flags().String("name", "", "Название клиента (обязательно)")
	customersCreateCmd.Flags().String("logo-url", "", "URL логотипа")
	customersCreateCmd.Flags().String("owner", "", "ID владельца")
	customersCreateCmd.Flags().String("status", "", "ID статуса")
	customersCreateCmd.Flags().String("tier", "", "ID уровня")

	customersUpdateCmd.Flags().String("name", "", "Новое название")
	customersUpdateCmd.Flags().String("logo-url", "", "Новый URL логотипа")
	customersUpdateCmd.Flags().String("owner", "", "Новый ID владельца")
	customersUpdateCmd.Flags().String("status", "", "Новый ID статуса")
	customersUpdateCmd.Flags().String("tier", "", "Новый ID уровня")

	customersNeedsCreateCmd.Flags().String("body", "", "Текст потребности (обязательно)")
	customersNeedsCreateCmd.Flags().String("customer", "", "ID клиента")
	customersNeedsCreateCmd.Flags().String("issue", "", "ID задачи")

	customersNeedsCmd.AddCommand(customersNeedsListCmd)
	customersNeedsCmd.AddCommand(customersNeedsCreateCmd)

	customersStatusesCmd.AddCommand(customersStatusesListCmd)
	customersTiersCmd.AddCommand(customersTiersListCmd)

	customersCmd.AddCommand(customersListCmd)
	customersCmd.AddCommand(customersViewCmd)
	customersCmd.AddCommand(customersCreateCmd)
	customersCmd.AddCommand(customersUpdateCmd)
	customersCmd.AddCommand(customersDeleteCmd)
	customersCmd.AddCommand(customersNeedsCmd)
	customersCmd.AddCommand(customersStatusesCmd)
	customersCmd.AddCommand(customersTiersCmd)
	rootCmd.AddCommand(customersCmd)
}
