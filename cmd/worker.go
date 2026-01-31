package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the worker",
	Long:  "Start the worker is a command that starts the worker",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Worker started")
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
