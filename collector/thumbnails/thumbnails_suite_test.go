package thumbnails_test

import (
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/mauricioabreu/video_samples/collector/filesystem"
	"github.com/mauricioabreu/video_samples/collector/thumbnails"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
)

func TestThumbnails(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Thumbnails Suite")
}

var _ = Describe("Thumbnails insert", func() {
	When("Adding to redis succeeds", func() {
		It("inserts the thumbnail", func() {
			file := &filesystem.File{
				Path:    "/thumbnails/bunny/0001.jpg",
				Dir:     "bunny",
				Data:    []byte("test_data"),
				ModTime: 1678103906,
			}
			uuid := func() string { return "1" }
			redisClient, redisMock := redismock.NewClientMock()
			redisMock.
				ExpectZAdd("thumbnails/bunny", redis.Z{Score: float64(file.ModTime), Member: "blob/1"}).
				SetVal(0)
			redisMock.
				ExpectSet("blob/1", []byte("test_data"), time.Duration(30)*time.Second).SetVal("OK")

			err := thumbnails.Insert(file, 30, uuid, redisClient)

			Expect(redisMock.ExpectationsWereMet()).To(Not(HaveOccurred()))
			Expect(err).To(Not(HaveOccurred()))
		})
	})
	When("Adding to set fails", func() {
		It("does not insert the thumbnail", func() {
			file := &filesystem.File{
				Path:    "/thumbnails/bunny/0001.jpg",
				Dir:     "bunny",
				Data:    []byte("test_data"),
				ModTime: 1678103906,
			}
			uuid := func() string { return "1" }
			redisClient, redisMock := redismock.NewClientMock()
			redisMock.
				ExpectZAdd("thumbnails/bunny", redis.Z{Score: float64(file.ModTime), Member: uuid()}).
				SetErr(errors.New("failed to execute zadd cmd"))

			err := thumbnails.Insert(file, 30, uuid, redisClient)
			Expect(err).To(HaveOccurred())
		})
	})
})
