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

func ExampleStates_Incr() {
	states := New()
	key := "127.0.0.1"
	_ = states.Incr(key)
	if v, ok := states.Counts.Load(key); ok && v == 1 {
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

func ExampleStates_IncrTime() {
	states := New()
	key := "127.0.0.1"
	_ = states.IncrTime(key, time.Now())
	if v, ok := states.Counts.Load(key); ok && v == 1 {
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

func ExampleStates_Decr() {
	now := time.Now()
	then := now.Add(-2 * time.Second)
	states := New()
	key := "127.0.0.1"
	_ = states.IncrTime("127.0.0.1", then)
	count, _ := states.Decr(1)
	if v, ok := states.Counts.Load(key); ok && v == 0 && count == 1 {
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
	if states.Counts.Size() != 0 || states.Seens.Size() != 0 && count != 1 {
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
		states.IncrTime(fmt.Sprintf("%d", rand.Int63()), now.Add(time.Duration(rand.Intn(100))*time.Second))
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
	if v, ok := states2.Counts.Load(key); !ok && v != 1 {
		t.Error("Expecting states2.Counts[key] == 1")
	}
}
