package consistenthashing

import (
	"testing"
)

func TestConsistentHashing(t *testing.T) {
	hr := NewHashRing(3)

	// Add nodes to the hash ring.
	hr.AddNode("Node1", 1)
	hr.AddNode("Node2", 1)
	hr.AddNode("Node3", 1)

	// Test key mappings.
	keys := []string{"key1", "key2", "key3"}
	expectedMappings := map[string]string{
		"key1": "Node3", // Replace with expected results for your hash function.
		"key2": "Node2",
		"key3": "Node2",
	}
	for _, key := range keys {
		node, _ := hr.GetNode(key)
		if node != expectedMappings[key] {
			t.Errorf("Key %s expected to map to %s but got %s", key, expectedMappings[key], node)
		}
	}

	// Test after removing a node.
	hr.RemoveNode("Node3")
	_, exists := hr.GetNode("key1")
	if !exists {
		t.Errorf("Expected a valid node for key1 after removing Node2")
	}
}

func TestEmptyRing(t *testing.T) {
	hr := NewHashRing(3)

	_, exists := hr.GetNode("key1")
	if exists {
		t.Errorf("Expected no node for key1 on an empty ring")
	}
}

func TestThreadSafety(t *testing.T) {
	hr := NewHashRing(3)

	go func() {
		hr.AddNode("Node1", 1)
	}()
	go func() {
		hr.RemoveNode("Node1")
	}()
	go func() {
		hr.GetNode("key1")
	}()
}
