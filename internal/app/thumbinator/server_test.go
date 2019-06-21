package thumbinator

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type dummyStore struct{}

func (ds dummyStore) GetThumb(streamName string) string {
	return "thumb_blob_here"
}

func (ds dummyStore) SaveThumb(stream Stream, timestamp int64, blob []byte) error {
	return nil
}

func TestRetrieveSnapshot(t *testing.T) {
	server := Server{store: dummyStore{}}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.showSnapshot)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Returned wrong status code: got %v wanted %v",
			status, http.StatusOK)
	}

	if body := rr.Body; body.String() != "thumb_blob_here" {
		t.Errorf("Returnted wrong body: got %v wanted %v", body.String(), "thumb_blob_here")
	}

	if headers := rr.Header(); headers.Get("Content-Type") != "image/jpeg" {
		t.Errorf("Returned wrong header: got %v wanted image/jpeg", headers.Get("Content-Type"))
	}
}
