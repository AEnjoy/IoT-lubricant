package constant

const (
	Etcd_DatastorePrefixConsumerReg     = "/" + AppName + "/consumers/"
	Etcd_DatastoreLeaderElectionPath    = "/" + AppName + "/leader_election"
	Etcd_DatastoreAssignmentPrefix      = "/" + AppName + "/assignments/"
	Etcd_DatastoreVirtualPartitionsConf = "/" + AppName + "/virtual_partitions"
)
