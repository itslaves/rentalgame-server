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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&env, "env", "E", "develop", "API server environment")
	viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env"))

	rootCmd.AddCommand(debug.Command())
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
