package thumbinator

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
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
	qb := parseQuery(r.URL)
	if qb.streamName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if qb.timestamp > 0 {
		thumb = s.store.GetThumbByTimestamp(qb.streamName, qb.timestamp)
	} else {
		thumb = s.store.GetThumb(qb.streamName)
	}
	w.Header().Add("Content-Length", strconv.Itoa(len(thumb)))
	w.Header().Add("Content-Type", "image/jpeg")
	fmt.Fprint(w, thumb)
}

// Serve start HTTP server to show thumbs
func Serve() {
	server := Server{store: newRedisStore()}
	http.HandleFunc("/", server.showSnapshot)
	log.Fatal(http.ListenAndServe(":8181", nil))
}
