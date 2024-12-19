package cluster

import (
	"fmt"
	"sync"

	m "github.com/mdkhanga/kvstore/models"
)

type cluster struct {
	mu         sync.Mutex
	clusterMap map[string]m.ClusterMember
}

type ICluster interface {
	AddToCluster(Hostname string, port int16) error
	RemoveFromCluster(Hostname string, port int16) error
	ListCluster() ([]m.ClusterMember, error)
}

func (c *cluster) AddToCluster(hostname string, port int16) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s:%d", hostname, port)
	if _, exists := c.clusterMap[key]; exists {
		return fmt.Errorf("member %s already exists in the cluster", key)
	}

	c.clusterMap[key] = m.ClusterMember{Host: hostname, Port: port}
	return nil
}

func (c *cluster) RemoveFromCluster(Hostname string, port int32) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s:%d", Hostname, port)
	if _, exists := c.clusterMap[key]; !exists {
		return fmt.Errorf("member %s does not exist in the cluster", key)
	}

	delete(c.clusterMap, key)
	return nil
}

func (c *cluster) ListCluster() ([]m.ClusterMember, error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	members := make([]m.ClusterMember, 0, len(c.clusterMap))
	for _, member := range c.clusterMap {
		members = append(members, member)
	}
	return members, nil
}

func New() *cluster {
	return &cluster{}
}
