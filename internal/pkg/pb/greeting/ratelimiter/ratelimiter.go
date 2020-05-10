package ratelimiter

import (
	"time"

	"github.com/juju/ratelimit"
)

type rateLimiter struct {
	interval time.Duration
	timeout  time.Duration
	capacity int64
	quantum  int64
	buckets  map[string]*ratelimit.Bucket
}

func New(
	interval time.Duration,
	timeout time.Duration,
	capacity int64,
	quantum int64,
) *rateLimiter {
	return &rateLimiter{
		interval: interval,
		timeout:  timeout,
		capacity: capacity,
		quantum:  quantum,
		buckets:  make(map[string]*ratelimit.Bucket),
	}
}

func (rl *rateLimiter) Wait(user string) bool {
	bucket, ok := rl.buckets[user]
	if !ok {
		bucket = ratelimit.NewBucketWithQuantum(rl.interval, rl.capacity, rl.quantum)
		rl.buckets[user] = bucket
	}

	return bucket.WaitMaxDuration(int64(1), rl.timeout)
}
