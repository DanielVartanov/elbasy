package main

import (
	"log"
	"time"
)

type throttler struct {
	BucketSize int

	leakyBucketQuota chan int8
}

func NewThrottler(bucketSize int) *throttler{
	instance := new(throttler)
	instance.BucketSize = bucketSize
	instance.leakyBucketQuota = make(chan int8, instance.channelSize())
	instance.initialise()
	return instance
}

func (self *throttler) Throttle(action func()) {
	log.Println("[Throttler] Waiting for quota")
	<- self.leakyBucketQuota
	log.Println("[Throttler] Quota is given. Performing an action")
	action()
	log.Println("[Throttler] Action is done")
}

func (self *throttler) initialise() {
	for i := 1; i <= self.channelSize(); i++ {
		self.leakyBucketQuota <- 1
	}

	go func() {
		for {
			self.leakyBucketQuota <- 1
			log.Println("[Throttler] Time has passed, released a quota")
			time.Sleep(1 * time.Second)
		}
	}()
}

func (self *throttler) channelSize() int {
/*
Why the channel is of size of LEAKY_BUCKET_SIZE - 1?
If the bucket is empty, then the emptier is ready to add one piece of quota any time
So in case of a burst of requests of the size of LEAKY_BUCKET_SIZE, the bucket will be filled, _immediately_ the emptier will add 1 quota and all requests of the count of LEAKY_BUCKET_SIZE will pass through, but the very next one will wait for the cooldown.
*/

	return self.BucketSize - 1
}
