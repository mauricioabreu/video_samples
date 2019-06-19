package cmd

import (
	"github.com/mauricioabreu/thumbinator/internal/app/thumbinator"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run HTTP server to retrieve your thumbs",
	Run: func(cmd *cobra.Command, args []string) {
		thumbinator.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
