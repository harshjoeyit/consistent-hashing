package consistenthashing

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type HashRing struct {
	nodes        map[uint32]string // Maps hash values to node addresses.
	sortedHashes []uint32          // Sorted hashes of all nodes.
	replicas     uint              // Number of virtual nodes per real node.
	mutex        sync.RWMutex      // To make it safe for concurrent use.
}

func NewHashRing(replicas uint) *HashRing {
	return &HashRing{
		nodes:        map[uint32]string{},
		sortedHashes: []uint32{},
		replicas:     replicas,
		mutex:        sync.RWMutex{},
	}
}

// AddNode adds a node to the hash ring
// The default weight is 1, the number of virtual nodes is weight * Virtualnodes
// Weight denotes the capacity of the node
func (hr *HashRing) AddNode(node string, weight float32) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	virtualnodes := int(weight * float32(hr.replicas))
	for i := 0; i < virtualnodes; i++ {
		// virtual node name is "node#i"
		virtualNode := fmt.Sprintf("%s#%d", node, i)
		h := hash(virtualNode)
		hr.nodes[h] = virtualNode
		hr.sortedHashes = append(hr.sortedHashes, h)
	}

	sort.Slice(hr.sortedHashes, func(i, j int) bool {
		return hr.sortedHashes[i] < hr.sortedHashes[j]
	})
}

// RemoveNode removes a node from the hash ring
func (hr *HashRing) RemoveNode(node string) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	// delete the node and its virtual nodes from map
	for i := 0; i < int(hr.replicas); i++ {
		virtualNode := fmt.Sprintf("%s#%d", node, i)
		delete(hr.nodes, hash(virtualNode))
	}

	updatedsortedHashes := make([]uint32, 0)
	for _, h := range hr.sortedHashes {
		if _, exists := hr.nodes[h]; exists {
			updatedsortedHashes = append(updatedsortedHashes, h)
		}
	}

	hr.sortedHashes = updatedsortedHashes

	sort.Slice(hr.sortedHashes, func(i, j int) bool {
		return hr.sortedHashes[i] < hr.sortedHashes[j]
	})
}

// GetNode returns the node for the given key
func (hr *HashRing) GetNode(key string) (string, bool) {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	if len(hr.sortedHashes) == 0 {
		return "", false
	}

	h := hash(key)
	// search for the first node hash that is greater than the key hash
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
		return strings.Join(parts[:len(parts)-1], "#")
	}

	return node
}

// hash creatse a new sha1 hash and convert to int32
func hash(s string) uint32 {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return (uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])
}
