package thumbinator

import (
	"fmt"
	"net/http"
)

func showSnapshot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "snapshot")
}
