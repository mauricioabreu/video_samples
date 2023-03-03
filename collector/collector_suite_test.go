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
			It("matches the extension", func() {
				Expect(watcher.MatchExt("/thumbnails/bunny/0001.jpg", []string{"jpg", "jpeg", "png"})).To(BeTrue())
			})
		})
		When("patterns list does not contain the extension", func() {
			It("does not match the extension", func() {
				Expect(watcher.MatchExt("/thumbnails/bunny/0001.jpg", []string{"png", "bmp"})).To(BeFalse())
			})
		})
	})
})
