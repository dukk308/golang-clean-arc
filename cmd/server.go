package cmd

import (
	"os"

	"github.com/dukk308/golang-clean-arc/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  "Start the server is a command that starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Setenv("TZ", "UTC"); err != nil {
			panic(err)
		}

		app := server.Bootstrap()
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
