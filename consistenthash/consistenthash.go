package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// HashRing represents the consistent hash ring.
type HashRing struct {
	nodes        []string
	replicas     int
	hashMap      map[uint32]string
	sortedHashes []uint32
}

// NewHashRing creates a new HashRing.
func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[uint32]string),
	}
}

// AddNode adds a node to the hash ring.
func (hr *HashRing) AddNode(node string) {
	for i := 0; i < hr.replicas; i++ {
		replicaKey := strconv.Itoa(i) + node
		hash := crc32.ChecksumIEEE([]byte(replicaKey))
		hr.nodes = append(hr.nodes, node)
		hr.hashMap[hash] = node
		hr.sortedHashes = append(hr.sortedHashes, hash)
	}

	sort.Slice(hr.sortedHashes, func(i, j int) bool {
		return hr.sortedHashes[i] < hr.sortedHashes[j]
	})
}

// RemoveNode removes a node from the hash ring.
func (hr *HashRing) RemoveNode(node string) {
	for i := 0; i < hr.replicas; i++ {
		replicaKey := strconv.Itoa(i) + node
		hash := crc32.ChecksumIEEE([]byte(replicaKey))

		// Find and remove the hash entry
		for j, h := range hr.sortedHashes {
			if h == hash {
				hr.sortedHashes = append(hr.sortedHashes[:j], hr.sortedHashes[j+1:]...)
				break
			}
		}

		// Remove the node entry
		delete(hr.hashMap, hash)
	}

	// Remove the node from the list
	for i, n := range hr.nodes {
		if n == node {
			hr.nodes = append(hr.nodes[:i], hr.nodes[i+1:]...)
			break
		}
	}
}

// GetNode returns the node to which the key is mapped.
func (hr *HashRing) GetNode(key string) string {
	if len(hr.nodes) == 0 {
		return ""
	}

	hash := crc32.ChecksumIEEE([]byte(key))

	// Binary search for the node with the smallest hash greater than or equal to the key's hash
	idx := sort.Search(len(hr.sortedHashes), func(i int) bool {
		return hr.sortedHashes[i] >= hash
	})

	// Wrap around if needed
	if idx == len(hr.sortedHashes) {
		idx = 0
	}

	return hr.hashMap[hr.sortedHashes[idx]]
}

/* func main() {
	// Create a new HashRing with 3 replicas
	hr := NewHashRing(3)

	// Add nodes to the hash ring
	hr.AddNode("server1")
	hr.AddNode("server2")
	hr.AddNode("server3")

	// Get node for a key
	key := "some_key"
	node := hr.GetNode(key)
	fmt.Printf("Key '%s' is mapped to node '%s'\n", key, node)

	// Remove a node from the hash ring
	nodeToRemove := "server2"
	hr.RemoveNode(nodeToRemove)

	// Get node for the same key after removing a node
	nodeAfterRemoval := hr.GetNode(key)
	fmt.Printf("After removing node '%s', key '%s' is mapped to node '%s'\n", nodeToRemove, key, nodeAfterRemoval)
} */
