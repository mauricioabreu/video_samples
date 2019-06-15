package main

import (
	"github.com/mauricioabreu/thumbinator/internal/app/thumbinator"
	log "github.com/sirupsen/logrus"
)

func main() {
	thumbsPath := "thumbnails"
	streams, err := thumbinator.GetStreams(thumbinator.JSONSource{File: "streams.json"})
	if err != nil {
		log.Fatalf("Could not retrieve streams to process: %s", err)
	}
	log.Infof("Retrieved the following streams: %+v", streams)
	for _, s := range streams {
		log.Debugf("Generating thumbs for %s with URL %s", s.Name, s.URL)
		thumbinator.GenerateThumb(s.URL, s.Name, thumbsPath)
	}
	thumbinator.CollectThumbs(thumbsPath)
}
