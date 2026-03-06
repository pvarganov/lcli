package cmd

import (
	"github.com/pavelvarganov/lcli/internal/config"
	"github.com/spf13/cobra"
)

var token string

var rootCmd = &cobra.Command{
	Use:   "lcli",
	Short: "CLI-утилита для работы с Linear",
	Long:  "lcli — полнофункциональная CLI-утилита для Linear через GraphQL API",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Linear API токен (или установите LINEAR_API_KEY)")
}

// GetToken возвращает токен: сначала флаг --token, затем env LINEAR_API_KEY и файл конфига.
func GetToken() string {
	if token != "" {
		return token
	}
	t, _ := config.LoadToken()
	return t
}
