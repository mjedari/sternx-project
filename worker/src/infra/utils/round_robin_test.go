package utils

import (
	"reflect"
	"sync"
	"testing"
)

func TestRoundRobin_Next(t *testing.T) {
	// arrange
	numbers := 6
	data := makePartitionSlice(numbers) // assuming makePartitionSlice returns a slice of int

	rr := NewRoundRobin(data)

	// act
	result := make([]int, 0, numbers)
	for i := 0; i < numbers; i++ {
		result = append(result, rr.Next())
	}

	// assert
	if !reflect.DeepEqual(data, result) {
		t.Errorf("Expected %#v, got %#v", data, result)
	}
}

// This is the helper function for executing the concurrent operations.
func executeRoundRobin(numbers int) []int {
	data := makePartitionSlice(numbers)

	rr := NewRoundRobin(data)

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(numbers)
	var results []int
	for i := 0; i < numbers; i++ {
		go func() {
			defer wg.Done()
			val := rr.Next()
			mu.Lock()
			results = append(results, val)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return results
}

func TestRoundRobin_Concurrency_DeepEqual(t *testing.T) {
	numbers := 10
	data := makePartitionSlice(numbers)
	results := executeRoundRobin(numbers)

	// assert
	if !reflect.DeepEqual(data, results) {
		t.Errorf("Expected %#v, got %#v", data, results)
	}
}

func TestRoundRobin_Concurrency_EqualIgnoreOrder(t *testing.T) {
	numbers := 10
	data := makePartitionSlice(numbers)
	results := executeRoundRobin(numbers)

	// assert
	if !equalIgnoreOrder(data, results) {
		t.Errorf("Expected %#v, got %#v", data, results)
	}
}

func equalIgnoreOrder(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	count := make(map[int]int)

	for _, v := range a {
		count[v]++
	}

	for _, v := range b {
		count[v]--
		if count[v] < 0 {
			return false
		}
	}

	return true
}
