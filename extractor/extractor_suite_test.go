package extractor_test

import (
	"testing"

	"github.com/mauricioabreu/video_samples/extractor"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExtractor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extractor Suite")
}

var _ = Describe("Extract resources from video", func() {
	Describe("Extract thumbs", func() {
		It("Start to extract thumbs", func() {
			opts := extractor.ThumbOptions{
				Input:   "http://localhost:8080/big_buck_bunny/playlist.m3u8",
				Scale:   "-1:360",
				Quality: 5,
			}
			runner := func(extractor.Command) error {
				return nil
			}
			err := extractor.ExtractThumbs("big_buck_bunny", opts, runner)
			Expect(err).To(Not(HaveOccurred()))
		})
	})
})
