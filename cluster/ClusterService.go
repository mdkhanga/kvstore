package cluster

import (
	m "github.com/mdkhanga/kvstore/models"
)

type cluster struct {
	clusterMap map[string]m.ClusterMember
}

type ICluster interface {
	AddToCluster(Hostname string, port int32) error
	RemoveFromCluster(Hostname string, port int32) error
	ListCluster() ([]m.ClusterMember, error)
}

func (c *cluster) AddToCluster(Hostname string, port int32) error {

	return nil
}

func (c *cluster) RemoveFromCluster(Hostname string, port int32) error {

	return nil
}

func (c *cluster) ListCluster() ([]m.ClusterMember, error) {

	var members []m.ClusterMember

	return members, nil
}

func New() *cluster {
	return &cluster{}
}
