package debug

import (
	"github.com/itslaves/rentalgames-server/route"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug API server",
		Long:  "Debug API server",
		Run: func(cmd *cobra.Command, args []string) {
			route.Route()
		},
	}

	cmd.Flags().IntP("port", "p", 7777, "A port of API server")
	viper.BindPFlag("port", cmd.Flags().Lookup("port"))

	return cmd
}
