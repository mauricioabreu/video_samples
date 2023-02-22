package cmd

import (
	"github.com/mauricioabreu/video_samples/tasks"
	"github.com/spf13/cobra"
)

func Work() *cobra.Command {
	return &cobra.Command{
		Use:   "work",
		Short: "Start workers to extract and process video samples",
		Run: func(cmd *cobra.Command, args []string) {
			tasks.StartWorker()
		},
	}
}
