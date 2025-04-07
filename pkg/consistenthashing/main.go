package consistenthashing

import (
	"crypto/sha1"
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"
	"sync"
)

type HashRing struct {
	// nodes Maps hash values to node addresses.
	//
	// To avoid collissions - use a Larger Hash Function: Consider using sha256 or sha512
	// and taking more bytes from the hash output. For example, you could use a uint64
	// instead of a uint32.
	nodes map[uint32]string

	// sortedHashes stores Sorted hashes of all nodes.
	sortedHashes []uint32

	// replicas denotes number of virtual nodes per real node.
	replicas uint

	// weights stores the weight of each node.
	weights map[string]float64

	// mutex makes it safe for concurrent use.
	mutex sync.RWMutex
}

func NewHashRing(replicas uint) *HashRing {
	return &HashRing{
		nodes:        map[uint32]string{},
		sortedHashes: []uint32{},
		replicas:     replicas,
		weights:      map[string]float64{},
		mutex:        sync.RWMutex{},
	}
}

// AddNode adds a node to the hash ring
// The default weight is 1, the number of virtual nodes is weight * Virtualnodes
// Weight denotes the capacity of the node
func (hr *HashRing) AddNode(node string, weight float64) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	virtualnodes := int(math.Round(weight * float64(hr.replicas)))

	for i := range virtualnodes {
		// virtual node name is "node#i"
		virtualNode := fmt.Sprintf("%s#%d", node, i)

		// Note: Collissions are possible (increasing virtual nodes can still
		// smooth out the distribution)
		h := hash(virtualNode)

		// add this virtual node to the hash ring
		hr.sortedHashes = append(hr.sortedHashes, h)
		hr.nodes[h] = virtualNode
	}

	hr.weights[node] = weight

	slices.Sort(hr.sortedHashes)
}

// RemoveNode removes a node from the hash ring
func (hr *HashRing) RemoveNode(node string) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	weight := hr.weights[node]

	// delete the node and its virtual nodes from map
	virtualnodes := int(math.Round(weight * float64(hr.replicas)))

	for i := range virtualnodes {
		virtualNode := fmt.Sprintf("%s#%d", node, i)
		h := hash(virtualNode)
		// delete the virtual node from the hash ring
		delete(hr.nodes, h)
	}

	updatedSortedHashes := make([]uint32, 0)
	for _, h := range hr.sortedHashes {
		if _, exists := hr.nodes[h]; exists {
			updatedSortedHashes = append(updatedSortedHashes, h)
		}
	}

	// update the sorted hashes
	hr.sortedHashes = updatedSortedHashes

	// delete the node from the weights map
	delete(hr.weights, node)

	slices.Sort(hr.sortedHashes)
}

// GetNode returns the node for the given key
func (hr *HashRing) GetNode(key string) (string, bool) {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	if len(hr.sortedHashes) == 0 {
		return "", false
	}

	h := hash(key)

	// search for the first node hash that is greater than the key hash 'h'
	i := sort.Search(len(hr.sortedHashes), func(i int) bool {
		return hr.sortedHashes[i] >= h
	})

	if i == len(hr.sortedHashes) {
		i = 0
	}

	node := hr.nodes[hr.sortedHashes[i]]

	return hr.getActualNode(node), true
}

// getActualNode returns the actual node name if the node is virtual node
// else returns the same node
func (*HashRing) getActualNode(node string) string {
	if strings.Contains(node, "#") {
		parts := strings.Split(node, "#")
		return parts[0]
	}

	return node
}

// hash returns a new integer sum for a sha1 hash
func hash(s string) uint32 {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return (uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])
}
