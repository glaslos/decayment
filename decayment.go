/*
Copyright 2016 Lukas Rist. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

// Package decayment decays counters after not being updated during threshold
package decayment

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/puzpuzpuz/xsync/v2"
)

// States to be decayed
type States struct {
	Counts     *xsync.MapOf[string, int]
	Seens      *xsync.MapOf[string, time.Time]
	tickerChan chan struct{}
}

// New state instance
func New() *States {
	s := States{
		Counts:     xsync.NewMapOf[int](),
		Seens:      xsync.NewMapOf[time.Time](),
		tickerChan: make(chan struct{}),
	}
	return &s
}

// Incr increments parameter key by one setting seen to now
func (s *States) Incr(key string) error {
	return s.IncrTime(key, time.Now())
}

// IncrTime increments parameter key by one setting seen to parameter t
// Used for testing to increment a key in the past
func (s *States) IncrTime(key string, t time.Time) error {
	s.Counts.Compute(key, func(oldValue int, loaded bool) (newValue int, delete bool) {
		// loaded is true here.
		newValue = oldValue + 1
		delete = false
		return
	})
	s.Seens.Store(key, t)
	return nil
}

// Decr decrements all keys if seen below threshold*time.Second
func (s *States) Decr(threshold int) (uint32, error) {
	count := new(uint32)
	s.Seens.Range(func(key string, seen time.Time) bool {
		interval := (time.Duration(threshold) * time.Second).Seconds()
		seenSince := int(time.Since(seen).Seconds() / interval)
		if seenSince >= 1 {
			s.Counts.Compute(key, func(oldValue int, loaded bool) (newValue int, delete bool) {
				// loaded is true here.
				delete = false
				newValue = oldValue - seenSince
				if newValue <= 0 {
					*count++
					s.Seens.Delete(key)
					delete = true
				}
				return
			})
		}
		return false
	})
	return *count, nil
}

// Start starts the decrementing loop with interval*time.Second
// and threshold*time.Second
func (s *States) Start(interval int, threshold int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	s.Decr(threshold)
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

// Encode the state as a byte array
func (s *States) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode the state from a byte array
func (s *States) Decode(b []byte) error {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&s)
}
