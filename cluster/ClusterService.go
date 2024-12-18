package cluster

import (
	m "github.com/mdkhanga/kvstore/models"
)

var clusterMap map[string]m.ClusterMember

func AddToCluster(Hostname string, port int32) error {

	return nil
}

func RemoveFromCluster(Hostname string, port int32) error {

	return nil
}

func ListCluster() ([]m.ClusterMember, error) {

	var members []m.ClusterMember

	return members, nil
}
