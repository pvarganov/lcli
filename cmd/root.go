package cmd

import (
	"github.com/pavelvarganov/lcli/internal/config"
	"github.com/spf13/cobra"
)

var token string
var outputFormat string

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
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Формат вывода: table (default), json")
}

// GetToken возвращает токен: сначала флаг --token, затем env LINEAR_API_KEY и файл конфига.
func GetToken() (string, error) {
	if token != "" {
		return token, nil
	}
	return config.LoadToken()
}

// GetOutputFormat возвращает текущий формат вывода.
func GetOutputFormat() string {
	return outputFormat
}
