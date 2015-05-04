package keen_test

import (
	"os"

	"github.com/oreillymedia/go-keen"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analysis", func() {
	var client *keen.Client
	var dest keen.AnalysisResult

	BeforeEach(func() {
		client = &keen.Client{
			ReadKey:   os.Getenv("KEEN_READ_KEY"),
			ProjectID: os.Getenv("KEEN_PROJECT_ID"),
		}
	})

	Describe("Count", func() {
		It("should get a basic count", func() {
			client.Count(&keen.AnalysisParams{
				EventCollection: "Loaded a Page",
			}, &dest)

			Expect(dest.Result).ToNot(BeZero())
		})

		It("should handle groupBy", func() {
			client.Count(&keen.AnalysisParams{
				EventCollection: "Loaded a Page",
				GroupBy:         "path",
			}, &dest)

			Expect(true).To(BeFalse())
		})

		It("should handle timeframe strings", func() {
			client.Count(&keen.AnalysisParams{
				EventCollection: "Loaded a Page",
				Timeframe:       "this_7_days",
			}, &dest)

			Expect(true).To(BeFalse())
		})

		It("should handle Interval", func() {
			client.Count(&keen.AnalysisParams{
				EventCollection: "Loaded a Page",
				Timeframe:       "this_7_days",
				Interval:        "daily",
			}, &dest)
			Expect(true).To(BeFalse())
		})
	})
})
