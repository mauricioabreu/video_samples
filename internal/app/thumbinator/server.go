package thumbinator

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Server HTTP handler
type Server struct {
	store store
}

func (s Server) showSnapshot(w http.ResponseWriter, r *http.Request) {
	thumb := s.store.GetThumb("big_buck_bunny")
	w.Header().Add("Content-Type", "image/jpeg")
	fmt.Fprint(w, thumb)
}

// Serve start HTTP server to show thumbs
func Serve() {
	server := Server{store: newRedisStore()}
	http.HandleFunc("/", server.showSnapshot)
	log.Fatal(http.ListenAndServe(":8181", nil))
}
