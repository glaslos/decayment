package decayment

import "fmt"

func ExampleIncrement() {
	states := New()
	key := "127.0.0.1"
	err := states.Incr(key)
	if err == nil && states.Counts[key] == 1 {
		fmt.Println("olleh")
	}
	// Output: olleh
}
