package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run HTTP server to retrieve your thumbs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serving...")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
