package utils

import (
	"sync"
)

/*
	Note:
	This implementation id not order deterministic version in concurrent usage
	But it guarantee it's order functionality in single thread usage and non-ordered in concurrent
*/

type RoundRobin struct {
	data []int
	i    int
	mu   sync.Mutex
}

func NewRoundRobin(data []int) *RoundRobin {
	return &RoundRobin{data: data}
}

func (r *RoundRobin) Next() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.data) <= 1 {
		return 0
	}
	val := r.data[r.i]
	r.i = (r.i + 1) % len(r.data)
	return val
}

func makePartitionSlice(numbers int) []int {
	var partitions []int
	for i := 0; i < numbers; i++ {
		partitions = append(partitions, i)
	}
	return partitions
}
