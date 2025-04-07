package main

import (
	"fmt"

	"github.com/harshjoeyit/myconsitenthashing/pkg/consistenthashing"
)

func main() {
	fmt.Println("Hello, World!")
	hr := consistenthashing.NewHashRing(50)

	hr.AddNode("server1", 1)
	hr.AddNode("server2", 1)
	hr.AddNode("server3", 2)

	counter := make(map[string]int)

	for i := range 100 {
		node, exists := hr.GetNode(fmt.Sprintf("key%d", i))
		if !exists {
			fmt.Println("Node not found")
			continue
		}
		counter[node]++
	}

	for key, val := range counter {
		fmt.Println(key, val)
	}

	hr.AddNode("server4", 1)

	counter = make(map[string]int)

	for i := range 100 {
		node, exists := hr.GetNode(fmt.Sprintf("key%d", i))
		if !exists {
			fmt.Println("Node not found")
			continue
		}
		counter[node]++
	}

	for key, val := range counter {
		fmt.Println(key, val)
	}

	hr.RemoveNode("server1")
	hr.RemoveNode("server2")
	hr.RemoveNode("server3")
	hr.RemoveNode("server4")

	node, exists := hr.GetNode("key1")
	if !exists {
		fmt.Println("Node not found")
		return
	}

	fmt.Println("node found", node)
}
