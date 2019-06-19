package thumbinator

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func showSnapshot(w http.ResponseWriter, r *http.Request) {
	thumb := getThumb("big_buck_bunny")
	w.Header().Add("Content-Type", "image/jpeg")
	fmt.Fprint(w, thumb)
}

func getThumb(streamName string) string {
	client := newClient()
	keys, err := client.ZRevRange("thumbs/"+streamName, 0, 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	thumb, err := client.Get("thumbs/blob/" + keys[0]).Result()
	if err != nil {
		log.Fatal(err)
	}
	return thumb
}

// Serve start HTTP server to show thumbs
func Serve() {
	http.HandleFunc("/", showSnapshot)
	log.Fatal(http.ListenAndServe(":8181", nil))
}
