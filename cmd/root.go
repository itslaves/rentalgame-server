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
	viper.AutomaticEnv()
	viper.SetEnvPrefix("rg")

	env := "develop"
	if viper.IsSet("env") {
		env = viper.GetString("env")
	}

	viper.AddConfigPath(fmt.Sprintf("./config/%s", env))
	viper.SetConfigName("application")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read the config file")
		os.Exit(1)
	}

	rootCmd.AddCommand(debug.Command())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
