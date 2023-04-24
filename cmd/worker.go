package cmd

import (
	"github.com/mauricioabreu/video_samples/config"
	"github.com/mauricioabreu/video_samples/tasks"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Work() *cobra.Command {
	return &cobra.Command{
		Use:   "work",
		Short: "Start workers to extract and process video samples",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.GetConfig()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get config")
			}
			tasks.StartWorker(&cfg)
		},
	}
}
