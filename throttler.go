package main

import (
	"log"
	"time"
)

const SHOPIFY_LEAKY_BUCKET_SIZE = 40
const SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND = 2

type throttler struct {
	bucketSize int
	leakRatePerSecond int
	storeQuotas map[string]*storeQuota
}

func NewThrottler() *throttler{
	instance := new(throttler)
	instance.bucketSize = SHOPIFY_LEAKY_BUCKET_SIZE
	instance.leakRatePerSecond = SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND
	instance.storeQuotas = make(map[string]*storeQuota)
	return instance
}

type storeQuota struct {
	quotaTrackingChannel chan int8
}

func newStoreQuota(bucketSize int, leakRatePerSecon int) *storeQuota {
	quotaTrackingChannel := make(chan int8, bucketSize)

	go func() {
		for {
			for i := 1; i <= leakRatePerSecon; i++ {
				quotaTrackingChannel <- 1
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return &storeQuota{quotaTrackingChannel: quotaTrackingChannel}
}

func (self *throttler) Throttle(store string, action func()) {
	storeQuota := self.fetchOrInitStoreQuota(store)

	start := time.Now()
	<- storeQuota.quotaTrackingChannel
	log.Printf("[Throttler] Waiting for quota took %s", time.Since(start))
	action()
}

func (self *throttler) fetchOrInitStoreQuota(store string) *storeQuota {
	if self.storeQuotas[store] == nil {
		self.storeQuotas[store] = newStoreQuota(self.channelSize(), SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND)
	}

	return self.storeQuotas[store]
}

func (self *throttler) channelSize() int {
/*
Why the channel is of size of LEAKY_BUCKET_SIZE - 1?
If the bucket is empty, then the emptier is ready to add one piece of quota any time
So in case of a burst of requests of the size of LEAKY_BUCKET_SIZE, the bucket will be filled, _immediately_ the emptier will add 1 quota and all requests of the count of LEAKY_BUCKET_SIZE will pass through, but the very next one will wait for the cooldown.

It is now `- LEAK_RATE` rather than `- 1`, but you got the idea.
*/

	return self.bucketSize - SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND
}
