package cmd

import (
	"fmt"
	"os"

	"github.com/itslaves/rentalgames-server/cmd/debug"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	env string
)

var rootCmd = &cobra.Command{
	Use:   "rentalgames",
	Short: "RentalGames API server",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&env, "env", "E", "develop", "API server environment")

	// --- BEGIN COMMANDS --- //
	rootCmd.AddCommand(debug.Command())
	// --- END COMMANDS --- //

	initConfig()
}

func initConfig() {
	if env := rootCmd.Flag("env"); env != nil {
		if v := env.Value.String(); v != "" {
			viper.AddConfigPath(fmt.Sprintf("./config/%s", v))
			viper.SetConfigName("application")
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can not read the config file")
		os.Exit(1)
	}
}

// Execute CLI 명령어 처리 수행
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
