package gobreaker

import (
	"github.com/sony/gobreaker"
)

const ( // Defaults
	DefaultAnalysisMinRequests = 200
	DefaultShouldOpenRatio     = 0.5
)

func ErrorRatioStrategy(minReqs uint32, openRatio float64) func(counts gobreaker.Counts) bool {
	if minReqs <= 0 {
		minReqs = DefaultAnalysisMinRequests
	}
	if openRatio <= 0 || openRatio > 1 {
		openRatio = DefaultShouldOpenRatio
	}
	return func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= minReqs && failureRatio >= openRatio
	}
}
