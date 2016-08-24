package decayment

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	New()
}

func TestIncrement(t *testing.T) {
	states := New()
	err := states.Incr("127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}

func TestIncrementTime(t *testing.T) {
	states := New()
	err := states.IncrTime("127.0.0.1", time.Now())
	if err != nil {
		t.Error(err)
	}
}

func TestDecrement(t *testing.T) {
	states := New()
	err := states.Decr(1)
	if err != nil {
		t.Error(err)
	}
}

func TestTrueDecrement(t *testing.T) {
	now := time.Now()
	then := now.Add(-2 * time.Second)
	states := New()
	err := states.IncrTime("127.0.0.1", then)
	if err != nil {
		t.Error(err)
	}
	err = states.Decr(1)
	if err != nil {
		t.Error(err)
	}
	if len(states.Counts) != 0 || len(states.Seens) != 0 {
		t.Error("state not properly decremented")
	}
}
