package main

import (
	"net"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	New()
}

func TestIncrement(t *testing.T) {
	states := New()
	err := states.incr(net.ParseIP("127.0.0.1"))
	if err != nil {
		t.Error(err)
	}
}

func TestIncrementTime(t *testing.T) {
	states := New()
	err := states.incrTime(net.ParseIP("127.0.0.1"), time.Now())
	if err != nil {
		t.Error(err)
	}
}

func TestDecrement(t *testing.T) {
	states := New()
	err := states.decr()
	if err != nil {
		t.Error(err)
	}
}

func TestTrueDecrement(t *testing.T) {
	now := time.Now()
	then := now.Add(-61 * time.Second)
	states := New()
	err := states.incrTime(net.ParseIP("127.0.0.1"), then)
	if err != nil {
		t.Error(err)
	}
	err = states.decr()
	if err != nil {
		t.Error(err)
	}
	if len(states.counts) != 0 || len(states.seens) != 0 {
		t.Error("state not properly decremented")
	}
}
