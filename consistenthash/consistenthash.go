package consistenthash

import (
	"hash/fnv"
	"sync"
)

type Node struct {
	ID string
}

type VirtualNode struct {
	ID     int
	NodeID string
}

type ConsistentHash struct {
	nodes        map[string]Node
	virtualNodes []VirtualNode
	totalVNodes  int
	mu           sync.RWMutex
}

func NewConsistentHash(totalVNodes int) *ConsistentHash {
	ch := &ConsistentHash{
		nodes:        make(map[string]Node),
		virtualNodes: make([]VirtualNode, totalVNodes),
		totalVNodes:  totalVNodes,
	}

	for i := 0; i < totalVNodes; i++ {
		ch.virtualNodes[i] = VirtualNode{ID: i, NodeID: ""}
	}
	return ch
}

func (ch *ConsistentHash) AddNode(node Node) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.nodes[node.ID] = node
	ch.redistributeVNodes()
}

func (ch *ConsistentHash) RemoveNode(nodeID string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	delete(ch.nodes, nodeID)
	ch.redistributeVNodes()
}

func (ch *ConsistentHash) redistributeVNodes() {
	if len(ch.nodes) == 0 {
		for i := range ch.virtualNodes {
			ch.virtualNodes[i].NodeID = ""
		}
		return
	}

	nodeIDs := make([]string, 0, len(ch.nodes))
	for id := range ch.nodes {
		nodeIDs = append(nodeIDs, id)
	}

	for i := range ch.virtualNodes {
		nodeIndex := i % len(nodeIDs)
		ch.virtualNodes[i].NodeID = nodeIDs[nodeIndex]
	}
}

func (ch *ConsistentHash) GetNode(key string) string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if len(ch.nodes) == 0 {
		return ""
	}

	hash := ch.hash(key)
	idx := hash % ch.totalVNodes

	return ch.virtualNodes[idx].NodeID
}

func (ch *ConsistentHash) hash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}
