# Consistent Hashing Implementation

This project implements consistent hashing, a technique used to distribute data across a cluster of nodes in a way that minimizes disruption when nodes are added or removed.

## Key Features

*   **Node Addition/Removal:**  Adding or removing nodes only requires remapping a small portion of the data.
*   **Virtual Nodes:**  Supports virtual nodes to improve distribution and handle nodes with varying capacities.
*   **Configurable Replicas:**  The number of virtual nodes (replicas) per physical node is configurable.
*   **Thread-Safe:**  Uses a mutex to ensure thread safety for concurrent access.
*   **SHA-1 Hashing:** Employs SHA-1 for a simple implementation. Can use SHA-1 or SHA-512 for reduced collision probability.

## Usage

1.  **Create a Hash Ring:**

    ```go
    ring := consistenthashing.NewHashRing(100) // 100 virtual nodes per physical node
    ```

2.  **Add Nodes:**

    ```go
    ring.AddNode("node1", 1.0) // Weight of 1.0
    ring.AddNode("node2", 1.5) // Weight of 1.5 (more capacity)
    ring.AddNode("node3", 1.0)
    ```

3.  **Get Node for a Key:**

    ```go
    node, ok := ring.GetNode("some-data-key")
    if ok {
        fmt.Println("Key 'some-data-key' belongs to node:", node)
    } else {
        fmt.Println("No nodes in the hash ring.")
    }
    ```

4.  **Remove a Node:**

    ```go
    ring.RemoveNode("node2")
    ```

## Implementation Details

*   **Hash Function:** SHA-1 is used to generate a 20 byte (160-bit hash), which is then truncated to a `uint32` for use in the hash ring.
*   **Virtual Nodes:** Each physical node is represented by multiple virtual nodes on the hash ring. The number of virtual nodes is determined by the `replicas` parameter and the node's weight.
*   **Hash Ring:** The hash ring is implemented as a sorted slice of `uint32` hash values.
*   **Node Lookup:**  The `GetNode` function uses a binary search to find the node whose hash value is closest to the hash of the key.

## Considerations

*   **Hash Collisions:** While SHA-1 provides a large hash space, collisions are still possible. Increasing the number of virtual nodes can help to mitigate the impact of collisions.
*   **Weighting:** Node weights allow you to assign more capacity to some nodes than others.  A node with a higher weight will receive a proportionally larger share of the data.

## Future Enhancements

*   Implement more sophisticated collision detection and resolution strategies.