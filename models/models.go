package models

type KeyValue struct {
	Key   string
	Value string
}

type ClusterMember struct {
	Host string
	Port int16
}
