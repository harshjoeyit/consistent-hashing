package consistenthashing

import (
	"strconv"
	"sync"
	"testing"
)

func TestStress(t *testing.T) {
	hr := NewHashRing(50)
	nodeCount := 10000
	keyCount := 100000
	var wg sync.WaitGroup

	// Add nodes.
	for i := 0; i < nodeCount; i++ {
		hr.AddNode("Node"+strconv.Itoa(i), 1)
	}

	// Query keys in parallel.
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(start int) {
			defer wg.Done()
			for j := start; j < keyCount; j += 10 {
				_, exists := hr.GetNode("key" + strconv.Itoa(j))
				if !exists {
					t.Errorf("Key %d should map to a valid node", j)
				}
			}
		}(i)
	}
	wg.Wait()

	// Remove nodes and verify.
	for i := 0; i < nodeCount; i++ {
		hr.RemoveNode("Node" + strconv.Itoa(i))
	}
	if len(hr.sortedHashes) != 0 {
		t.Errorf("All nodes should have been removed")
	}
}
