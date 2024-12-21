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
	hr.AddNode("server3", 1)

	counter := make(map[string]int)

	for i := 0; i < 100; i++ {
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

	for i := 0; i < 100; i++ {
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
}
