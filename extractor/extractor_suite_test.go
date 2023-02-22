package extractor_test

import (
	"errors"
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
		When("Running it works", func() {
			It("Start to extract thumbs", func() {
				opts := extractor.ThumbOptions{
					Input:   "http://localhost:8080/colors/playlist.m3u8",
					Scale:   "-1:360",
					Quality: 5,
				}
				runner := func(extractor.Command) error {
					return nil
				}
				err := extractor.ExtractThumbs("colors", opts, runner)
				Expect(err).To(Not(HaveOccurred()))
			})
		})
		When("Running it fails", func() {
			It("Start to extract thumbs", func() {
				opts := extractor.ThumbOptions{
					Input:   "http://localhost:8080/colors/playlist.m3u8",
					Scale:   "-1:360",
					Quality: 5,
				}
				runner := func(extractor.Command) error {
					return errors.New("failed to run")
				}
				err := extractor.ExtractThumbs("colors", opts, runner)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
