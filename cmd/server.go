package cmd

import (
	"os"

	"github.com/dukk308/golang-clean-arch-starter/internal/server"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/utils"
	"github.com/spf13/cobra"

	_ "github.com/dukk308/golang-clean-arch-starter/api-docs/swagger" //nolint
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  "Start the server is a command that starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Setenv("TZ", "UTC"); err != nil {
			panic(err)
		}

		utils.ParseFlags()
		app := server.Bootstrap(cmd.Context())
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
