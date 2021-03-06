package ratelimit

import (
	"sync"
	"sync/atomic"
	"time"
)

type Item struct {
	num int64
	expireIn int64
}

type Ratelimiter struct {
	s sync.Map
	duration int64
	rateLimit int64
	gcInterval time.Duration
}

func NewLimiter(duration, rateLimit int64, gcInterval time.Duration) *Ratelimiter {
	s := &Ratelimiter{
		duration: duration,
		rateLimit: rateLimit,
		gcInterval: gcInterval,
	}

	go s.cleaner()
	return s
}

func (s *Ratelimiter) cleaner() {
	t := time.NewTicker(s.gcInterval)
	defer t.Stop()
	for {
		select {
			case <-t.C:
				s.s.Range(func(k, v interface{}) bool {
					vv, _ := v.(*Item)
					if vv.expireIn < time.Now().Unix() {
						s.s.Delete(k)
					}
					return true
				})
		}
	}
}

func (s *Ratelimiter) incr(k string) {
	v, ok := s.s.LoadOrStore(k, &Item{num: 1, expireIn: time.Now().Unix() + s.duration })
	if ok {
		vv, _ := v.(*Item)
		atomic.AddInt64(&vv.num, 1)
	}
}

func (s *Ratelimiter) ShouldLimit(k string) bool {
	i, ok := s.s.Load(k)
	if !ok {
		s.incr(k)
		return false
	}
	s.incr(k)
	iv, _ := i.(*Item)
	if iv.expireIn <= time.Now().Unix() {
		s.s.Delete(k)
		s.incr(k)
		return false
	}
	if atomic.LoadInt64(&iv.num) > s.rateLimit {
		return true
	}
	return false
}