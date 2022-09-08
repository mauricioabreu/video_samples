package cmd

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mauricioabreu/video_samples/config"
	"github.com/mauricioabreu/video_samples/internal/app/video_samples"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var thumbsPath string
var streamsFile string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate thumbs from live streamings and videos on demand",
	Run: func(_ *cobra.Command, _ []string) {
		Main()
	},
}

// Run video_samples, run
func Run() {
	if err := generateCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

// Main run video_samples, run
func Main() {
	c := config.GetConfig()
	streams, err := video_samples.GetStreams(video_samples.JSONSource{File: streamsFile})
	if err != nil {
		log.Fatal().Err(err).Msg("Could not retrieve streams to process")
	}
	log.Info().Msgf("Retrieved the following streams: %+v", streams)
	for _, s := range streams {
		log.Debug().Msgf("Generating thumbs for %s with URL %s", s.Name, s.URL)
		video_samples.GenerateThumb(s.URL, s.Name, thumbsPath)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	redis, err := video_samples.NewRedisStore(c)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	collector := video_samples.Collector{Watcher: watcher, Store: redis, Path: thumbsPath}
	collector.CollectThumbs(streams)
}

func init() {
	generateCmd.Flags().StringVar(&thumbsPath, "thumbsPath", "thumbnails", "Path where all thumbs will be written to")
	generateCmd.Flags().StringVar(&streamsFile, "streamsFile", "streams.json", "File with streams to extract thumbs")
	rootCmd.AddCommand(generateCmd)
}
