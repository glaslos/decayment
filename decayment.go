package main

import (
	"encoding/binary"
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
	counts map[uint32]int64
	seens  map[uint32]time.Time
}

func (s *States) incr(ip net.IP) error {
	s.Lock()
	defer s.Unlock()
	intIP := ip2int(ip)
	s.counts[intIP]++
	s.seens[intIP] = time.Now()
	return nil
}

func (s *States) decr() error {
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

func main() {
	states := new(States)
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				states.decr()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
