package decayment

import (
	"sync"
	"time"
)

// States engine
type States struct {
	sync.Mutex
	Counts     map[interface{}]int64
	Seens      map[interface{}]time.Time
	tickerChan chan struct{}
}

// New state
func New() *States {
	s := States{
		Counts:     make(map[interface{}]int64),
		Seens:      make(map[interface{}]time.Time),
		tickerChan: make(chan struct{}),
	}
	return &s
}

// Incr increments parameter key by one setting seen to now
func (s *States) Incr(key interface{}) error {
	return s.IncrTime(key, time.Now())
}

// IncrTime increments parameter key by one setting seen to parameter t
func (s *States) IncrTime(key interface{}, t time.Time) error {
	s.Lock()
	defer s.Unlock()
	s.Counts[key]++
	s.Seens[key] = t
	return nil
}

// Decr decrements all keys if seen below threshold*time.Second
func (s *States) Decr(threshold int) error {
	s.Lock()
	defer s.Unlock()
	for intIP, seen := range s.Seens {
		if time.Since(seen) >= time.Duration(threshold)*time.Second {
			s.Counts[intIP]--
			if s.Counts[intIP] <= 0 {
				delete(s.Counts, intIP)
				delete(s.Seens, intIP)
			}
		}
	}
	return nil
}

// Start starts the decrement loop with interval*time.Second and threshold*time.Second
func (s *States) Start(interval int, threshold int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.Decr(threshold)
			case <-s.tickerChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the decrement loop
func (s *States) Stop() {
	close(s.tickerChan)
}
