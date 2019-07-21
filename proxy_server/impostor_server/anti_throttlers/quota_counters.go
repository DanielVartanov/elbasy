package anti_throttlers

import (
	"time"
)

const SAFETY_FACTOR = 1.05 // Quite often API server miscalculate time just by a bit but enough to make you hit a wall of status code 429

type quotaCounter interface {
	ClaimQuota(quotaAccount string)
	ReleaseQuota(quotaAccount string)
}

type leakyBucket struct {
	Size int
	LeakRatePerSecond int

	quotaAccountChannels map[string]chan int8
}

func newLeakyBucket(size int, leakRatePerSecond int) quotaCounter {
	lk := leakyBucket{Size: size, LeakRatePerSecond: leakRatePerSecond}
	lk.quotaAccountChannels = make(map[string]chan int8)
	return lk
}

func (lk leakyBucket) ClaimQuota(quotaAccount string) {
	<- lk.quotaAccountChannel(quotaAccount)
}

func (lk leakyBucket) ReleaseQuota(quotaAccount string) {
	//do nothing
}

func (lk leakyBucket) quotaAccountChannel(quotaAccount string) chan int8 {
	if lk.quotaAccountChannels[quotaAccount] == nil {
		lk.quotaAccountChannels[quotaAccount] = make(chan int8, lk.quotaAccountChannelSize())

		go func() {
			for {
				for i := 1; i <= lk.LeakRatePerSecond; i++ {
					lk.quotaAccountChannels[quotaAccount] <- 1
				}
				time.Sleep(1000 * SAFETY_FACTOR * time.Millisecond)
			}
		}()
	}

	return lk.quotaAccountChannels[quotaAccount]
}

func (lk leakyBucket) quotaAccountChannelSize() int {
	/*
          Why the channel is of size of Size - 1?
          If the bucket is empty, then the emptier is ready to add one piece of quota any time
          So in case of a burst of requests of the size of Size, the bucket will be filled, _immediately_ the emptier will add 1 quota and all requests of the count of Size will pass through, but the very next one will wait for the cooldown.

          It is now `Size - LeakRate` rather than `Size - 1`, but you got the idea.
        */
	return lk.Size - lk.LeakRatePerSecond
}
