package cmd

import (
	"github.com/mauricioabreu/video_samples/collector"
	"github.com/spf13/cobra"
)

func Collect() *cobra.Command {
	return &cobra.Command{
		Use:   "collect",
		Short: "Watch, collect and store resources",
		Run: func(cmd *cobra.Command, args []string) {
			collector.Collect("testvideo/thumbs/colors")
		},
	}
}
