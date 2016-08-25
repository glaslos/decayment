package decayment

import (
	"fmt"
	"math/rand"
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

func TestStartStop(t *testing.T) {
	states := New()
	states.Start(1, 1)
	states.Stop()
}

func BenchmarkIncr(b *testing.B) {
	states := New()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		states.Incr("127.0.0.1")
	}
}

func BenchmarkDecr(b *testing.B) {
	states := New()
	now := time.Now()
	for n := 0; n < 1000; n++ {
		states.IncrTime(rand.Int63(), now.Add(time.Duration(rand.Intn(100))*time.Second))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		states.Decr(1)
	}
}

func TestEncode(t *testing.T) {
	states := New()
	states.Incr("127.0.0.1")
	_, err := states.Encode()
	if err != nil {
		t.Error(err)
	}
}

func TestDecode(t *testing.T) {
	key := "127.0.0.1"
	states1 := New()
	states1.Incr(key)
	b, err := states1.Encode()
	if err != nil {
		t.Error(err)
	}
	states2 := New()
	err = states2.Decode(b)
	if err != nil {
		t.Error(err)
	}
	if states2.Counts[key] != 1 {
		t.Error("Expecting states2.Counts[key] == 1")
	}
}
