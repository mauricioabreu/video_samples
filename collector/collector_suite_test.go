package collector_test

import (
	"testing"

	"github.com/mauricioabreu/video_samples/collector/watcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCollector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collector Suite")
}

var _ = Describe("Collector", func() {
	Describe("Match file extension", func() {
		When("patterns list contains the extension", func() {
			It("matches", func() {
				patterns := []string{"jpg", "jpeg", "png"}
				Expect(watcher.MatchExt("/thumbnails/bunny/0001.jpg", patterns)).To(BeTrue())
			})
		})
		When("patterns list does not contain the extension", func() {
			It("does not match", func() {
				patterns := []string{"png", "bmp"}
				Expect(watcher.MatchExt("/thumbnails/bunny/0001.jpg", patterns)).To(BeFalse())
			})
		})
		When("path does not have extension", func() {
			It("does not match", func() {
				patterns := []string{"jpg", "jpeg", "png"}
				Expect(watcher.MatchExt("/thumbnails/bunny/0001", patterns)).To(BeFalse())
			})
		})
	})
})
