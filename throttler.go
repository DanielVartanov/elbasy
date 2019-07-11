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
	quotaTrackingChannel chan int8
}

func NewThrottler() *throttler{
	instance := new(throttler)
	instance.bucketSize = SHOPIFY_LEAKY_BUCKET_SIZE
	instance.leakRatePerSecond = SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND
	instance.quotaTrackingChannel = make(chan int8, instance.channelSize())
	instance.initialise()
	return instance
}

func (self *throttler) Throttle(action func()) {
	start := time.Now()
	<- self.quotaTrackingChannel
	log.Printf("[Throttler] Waiting for quota took %s", time.Since(start))
	action()
}

func (self *throttler) initialise() {
	go func() {
		for {
			for i := 1; i <= SHOPIFY_LEAKY_BUCKET_LEAK_RATE_PER_SECOND; i++ {
				self.quotaTrackingChannel <- 1
			}
			time.Sleep(1 * time.Second)
		}
	}()
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
