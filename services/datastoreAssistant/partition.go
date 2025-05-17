package datastoreAssistant

import (
	"sync"

	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/bytedance/sonic"
)

const defaultVPartitions = 64

var (
	assignedPartitions     = make(map[int]struct{})
	assignedPartitionsLock sync.RWMutex
)

// isPartitionAssigned checks if a given partition ID is assigned to this consumer
func isPartitionAssigned(partitionID int) bool {
	assignedPartitionsLock.RLock()
	defer assignedPartitionsLock.RUnlock()
	_, exists := assignedPartitions[partitionID]
	return exists
}

// updateAssignedPartitions updates the global state based on data from etcd
func updateAssignedPartitions(jsonData []byte, consumerID string) {
	var newPartitionList []int
	if err := sonic.Unmarshal(jsonData, &newPartitionList); err != nil {
		logg.L.Infof("[%s] Failed to unmarshal assignment data '%s': %v", consumerID, string(jsonData), err)
		return
	}

	newAssignments := make(map[int]struct{}, len(newPartitionList))
	for _, p := range newPartitionList {
		newAssignments[p] = struct{}{}
	}

	assignedPartitionsLock.Lock()
	assignedPartitions = newAssignments
	assignedPartitionsLock.Unlock()

	logg.L.Infof("[%s] Updated assigned partitions. Count: %d", consumerID, len(newPartitionList))
}

// clearAssignedPartitions clears the assignment state
func clearAssignedPartitions(consumerID string) {
	assignedPartitionsLock.Lock()
	assignedPartitions = make(map[int]struct{})
	assignedPartitionsLock.Unlock()
	logg.L.Infof("[%s] Cleared assigned partitions.", consumerID)
}
