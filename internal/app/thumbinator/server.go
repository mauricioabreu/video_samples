package thumbinator

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Server HTTP handler
type Server struct {
	store store
}

func (s Server) showSnapshot(w http.ResponseWriter, r *http.Request) {
	var thumb string
	timestamp, ok := r.URL.Query()["timestamp"]
	if ok {
		n, _ := strconv.ParseInt(timestamp[0], 10, 64)
		thumb = s.store.GetThumbByTimestamp("big_buck_bunny", n)
	} else {
		thumb = s.store.GetThumb("big_buck_bunny")
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
