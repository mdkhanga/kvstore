package models

import (
	"net"
)

type KeyValue struct {
	Key   string
	Value string
}

type ClusterMember struct {
	host string
	port int16
	conn net.Conn
}

type Cluster struct {
	members []ClusterMember
}
