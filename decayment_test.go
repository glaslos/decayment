package decayment

import (
	"fmt"
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

func ExampleIncrement() {
	states := New()
	key := "127.0.0.1"
	err := states.Incr(key)
	if err == nil && states.Counts[key] == 1 {
		fmt.Println("olleh")
	}
	// Output: olleh
}

func TestIncrementTime(t *testing.T) {
	states := New()
	err := states.IncrTime("127.0.0.1", time.Now())
	if err != nil {
		t.Error(err)
	}
}

func ExampleIncrementTime() {
	states := New()
	key := "127.0.0.1"
	err := states.IncrTime(key, time.Now())
	if err == nil && states.Counts[key] == 1 {
		fmt.Println("olleh")
	}
	// Output: olleh
}

func TestDecrement(t *testing.T) {
	states := New()
	_, err := states.Decr(1)
	if err != nil {
		t.Error(err)
	}
}

func ExampleDecrement() {
	now := time.Now()
	then := now.Add(-2 * time.Second)
	states := New()
	key := "127.0.0.1"
	err := states.IncrTime("127.0.0.1", then)
	count, err := states.Decr(1)
	if err == nil && states.Counts[key] == 0 && count == 1 {
		fmt.Println("olleh")
	}
	// Output: olleh
}

func TestTrueDecrement(t *testing.T) {
	now := time.Now()
	then := now.Add(-2 * time.Second)
	states := New()
	err := states.IncrTime("127.0.0.1", then)
	if err != nil {
		t.Error(err)
	}
	count, err := states.Decr(1)
	if err != nil {
		t.Error(err)
	}
	if len(states.Counts) != 0 || len(states.Seens) != 0 && count != 1 {
		t.Error("state not properly decremented")
	}
}
