package video_samples

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mauricioabreu/video_samples/config"
	"github.com/rs/zerolog/log"
)

// Server HTTP handler
type Server struct {
	store store
}

// queryBy arguments used to retrieve thumbs
type queryBy struct {
	streamName string
	timestamp  int64
}

func parseQuery(u *url.URL) queryBy {
	qb := queryBy{}
	streamName := u.Query().Get("stream_name")
	qb.streamName = streamName
	timestamp := u.Query().Get("timestamp")
	if timestamp != "" {
		n, _ := strconv.ParseInt(timestamp, 10, 64)
		qb.timestamp = n
	}
	return qb
}

func (s Server) showSnapshot(w http.ResponseWriter, r *http.Request) {
	var thumb string
	var err error
	var maxAge int
	qb := parseQuery(r.URL)
	if qb.streamName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if qb.timestamp > 0 {
		thumb, err = s.store.GetThumbByTimestamp(qb.streamName, qb.timestamp)
		maxAge = 3600
	} else {
		thumb, err = s.store.GetThumb(qb.streamName)
		maxAge = 4
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "")
		return
	}

	w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
	w.Header().Add("Content-Length", strconv.Itoa(len(thumb)))
	w.Header().Add("Content-Type", "image/jpeg")
	w.Header().Add("Expires", time.Now().Add(time.Second*time.Duration(maxAge)).Format(http.TimeFormat))
	fmt.Fprint(w, thumb)
}

// Serve start HTTP server to show thumbs
func Serve() {
	c := config.GetConfig()
	redis, err := NewRedisStore(c)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	server := Server{store: redis}
	http.HandleFunc("/", server.showSnapshot)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
