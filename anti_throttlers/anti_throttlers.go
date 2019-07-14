package anti_throttlers

import (
	"log"
	"time"
)

type AntiThrottler interface {
	PreventThrottling(quotaAccount string, action func())
}

const SHOPIFY_LEAKY_BUCKET_SIZE = 40
const SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND = 2

func NewShopifyAntiThrottler() AntiThrottler{
	return antiThrottler{
		quotaCounter: newLeakyBucket(
			SHOPIFY_LEAKY_BUCKET_SIZE,
			SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND,
		),
	}
}

type antiThrottler struct {
	quotaCounter quotaCounter
}

func (ath antiThrottler) PreventThrottling(quotaAccount string, action func()) {
	logElapsedTimeAsMetric("quota.claim_waiting_time", func() {
		ath.quotaCounter.ClaimQuota(quotaAccount)
	})
	action()
	ath.quotaCounter.ReleaseQuota(quotaAccount)
}

func logElapsedTimeAsMetric(metric_name string, action func()) {
 	start := time.Now()
	action()
	elapsed := time.Since(start).Seconds()
	log.Printf("%s=%.9fs", metric_name, elapsed)
}
