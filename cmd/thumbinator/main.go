package cmd

import (
	"os"

	"github.com/mauricioabreu/thumbinator/internal/app/thumbinator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var thumbsPath string
var streamsFile string

var runCmd = &cobra.Command{
	Use:   "thumbinator",
	Short: "Generate thumbs from live streamings and videos on demand",
	Run: func(cmd *cobra.Command, args []string) {
		Main()
	},
}

// Run thumbinator, run
func Run() {
	runCmd.Flags().StringVar(&thumbsPath, "thumbsPath", "thumbnails", "Path where all thumbs will be written to")
	runCmd.Flags().StringVar(&streamsFile, "streamsFile", "streams.json", "File with streams to extract thumbs")
	if err := runCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// Main run thumbinator, run
func Main() {
	streams, err := thumbinator.GetStreams(thumbinator.JSONSource{File: streamsFile})
	if err != nil {
		log.Fatalf("Could not retrieve streams to process: %s", err)
	}
	log.Infof("Retrieved the following streams: %+v", streams)
	for _, s := range streams {
		log.Debugf("Generating thumbs for %s with URL %s", s.Name, s.URL)
		thumbinator.GenerateThumb(s.URL, s.Name, thumbsPath)
	}
	thumbinator.CollectThumbs(streams, thumbsPath)
}
