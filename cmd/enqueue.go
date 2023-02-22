package cmd

import (
	"github.com/mauricioabreu/video_samples/tasks"
	"github.com/spf13/cobra"
)

func EnqueueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enqueue",
		Short: "Enqueue tasks to extract, collect and store video samples",
		Run: func(cmd *cobra.Command, args []string) {
			tasks.Enqueue()
		},
	}
}
