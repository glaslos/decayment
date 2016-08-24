package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// States engine
type States struct {
	sync.Mutex
	counts     map[uint32]int64
	seens      map[uint32]time.Time
	tickerChan chan struct{}
}

// New state
func New() *States {
	s := States{
		counts:     make(map[uint32]int64),
		seens:      make(map[uint32]time.Time),
		tickerChan: make(chan struct{}),
	}
	return &s
}

func (s *States) incr(ip net.IP) error {
	return s.incrTime(ip, time.Now())
}

func (s *States) incrTime(ip net.IP, t time.Time) error {
	s.Lock()
	defer s.Unlock()
	intIP := ip2int(ip)
	s.counts[intIP]++
	s.seens[intIP] = t
	return nil
}

func (s *States) decr() error {
	s.Lock()
	defer s.Unlock()
	for intIP, seen := range s.seens {
		if time.Since(seen) >= 60*time.Second {
			s.counts[intIP]--
			if s.counts[intIP] <= 0 {
				delete(s.counts, intIP)
				delete(s.seens, intIP)
			}
		}
	}
	return nil
}

func (s *States) start() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.decr()
			case <-s.tickerChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *States) stop() {
	close(s.tickerChan)
}

func main() {
	states := New()
	states.start()
	defer states.stop()
	states.incr(net.ParseIP("127.0.0.1"))
	fmt.Println(states)
	time.Sleep(70 * time.Second)
	fmt.Println(states)
}
