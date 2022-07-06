package cmd

import (
	"github.com/mauricioabreu/video_samples/internal/app/video_samples"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run HTTP server to retrieve your thumbs",
	Run: func(cmd *cobra.Command, args []string) {
		video_samples.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
