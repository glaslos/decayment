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
	"sync"
	"time"
)

// States to be decayed
type States struct {
	lock sync.Mutex
	// Countes per key
	Counts map[interface{}]int64
	// Seens keeps last seen per key
	Seens      map[interface{}]time.Time
	tickerChan chan struct{}
}

// New state instance
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

// Lock the states
func (s *States) Lock() {
	s.lock.Lock()
}

// Unlock the states
func (s *States) Unlock() {
	s.lock.Unlock()
}

// IncrTime increments parameter key by one setting seen to parameter t
// Used for testing to increment a key in the past
func (s *States) IncrTime(key interface{}, t time.Time) error {
	s.Lock()
	defer s.Unlock()
	s.Counts[key]++
	s.Seens[key] = t
	return nil
}

// Decr decrements all keys if seen below threshold*time.Second
func (s *States) Decr(threshold int) (uint32, error) {
	s.Lock()
	defer s.Unlock()
	count := new(uint32)
	for intIP, seen := range s.Seens {
		if time.Since(seen) >= time.Duration(threshold)*time.Second {
			s.Counts[intIP]--
			if s.Counts[intIP] <= 0 {
				*count++
				delete(s.Counts, intIP)
				delete(s.Seens, intIP)
			}
		}
	}
	return *count, nil
}

// Start starts the decrementing loop with interval*time.Second
// and threshold*time.Second
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

// Encode the state as a byte array
func (s *States) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode the state from a byte array
func (s *States) Decode(b []byte) error {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&s)
	if err != nil {
		return err
	}
	return nil
}
