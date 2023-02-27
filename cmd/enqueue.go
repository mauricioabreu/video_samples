package cmd

import (
	"github.com/mauricioabreu/video_samples/collector"
	"github.com/mauricioabreu/video_samples/extractor/inventory"
	"github.com/mauricioabreu/video_samples/tasks"
	"github.com/spf13/cobra"
)

func EnqueueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enqueue",
		Short: "Enqueue tasks to extract, collect and store video samples",
		Run: func(cmd *cobra.Command, args []string) {
			getStreams := func(url string) func() ([]inventory.Streaming, error) {
				return func() ([]inventory.Streaming, error) {
					return inventory.GetStreams(url)
				}
			}
			tasks.Enqueue(getStreams("http://localhost:8080/output.m3u8"))
			collector.Collect("testvideo/thumbs/colors")
		},
	}
}
