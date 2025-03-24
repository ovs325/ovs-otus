package cli

import (
	"fmt"
	"os"

	cb "github.com/spf13/cobra"
	cf "main/config"
	"main/internal/cli/add"
	"main/internal/cli/del"
	"main/internal/cli/reset"
	pr "main/internal/params"
)

const configPath = "./"

var rootCmd = &cb.Command{
	Use:   "AntiBruteForce",
	Short: "АнтиПодбор Паролей",
	Long:  "Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе",
}

func Execute() {
	config, err := cf.LoadConfig(configPath)
	pr.ConfigHTTP = config.HTTPServer
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}
	rootCmd.AddCommand(reset.ResetCmd)
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(del.DelCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
