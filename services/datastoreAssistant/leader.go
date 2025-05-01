package datastoreAssistant

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/bytedance/sonic"
	clientv3 "go.etcd.io/etcd/client/v3"

	"go.etcd.io/etcd/client/v3/concurrency"
)

func (a *app) runLeaderTasks(ctx context.Context, consumerID string) {
	logg.L.Infof("[%s] Acquired leadership. Starting leader tasks...", consumerID)
	defer logg.L.Infof("[%s] Stopping leader tasks.", consumerID)

	// Perform initial assignment immediately upon becoming leader
	a.performReassignment(consumerID)

	// Watch for changes in consumer registration
	consumerWatchChan := a.cli.Watch(ctx, constant.Etcd_DatastorePrefixConsumerReg, clientv3.WithPrefix(), clientv3.WithCreatedNotify()) // Start watching *after* initial assignment

	for {
		select {
		case <-ctx.Done():
			logg.L.Infof("[%s] Leader context cancelled. Exiting leader loop.", consumerID)
			return
		case watchResp, ok := <-consumerWatchChan:
			if !ok {
				logg.L.Warnf("[%s] Consumer watch channel closed. Exiting leader loop.", consumerID)
				return // Channel closed, stop leader tasks
			}
			if err := watchResp.Err(); err != nil {
				logg.L.Errorf("[%s] Consumer watch error: %v. Exiting leader loop.", consumerID, err)
				// todo: consider restart logic or graceful shutdown
				return
			}

			hasChanges := false
			for _, event := range watchResp.Events {
				logg.L.Infof("[%s] Consumer watch event: Type[%s] Key[%s]", consumerID, event.Type, string(event.Kv.Key))
				if event.Type == clientv3.EventTypePut || event.Type == clientv3.EventTypeDelete {
					// A consumer was added or removed
					hasChanges = true
				}
			}

			if hasChanges {
				logg.L.Infof("[%s] Detected consumer changes, triggering reassignment.", consumerID)
				a.performReassignment(consumerID)
			}
		}
	}
}
func (a *app) campaignForLeadership(consumerID string) {
	logg.L.Infof("[%s] Starting leader election campaign...", consumerID)

	session, err := concurrency.NewSession(a.cli, concurrency.WithTTL(leaseTTL+5))
	if err != nil {
		logg.L.Errorf("[%s] Failed to create concurrency session: %v", consumerID, err)
		return
	}
	defer session.Close()
	logg.L.Infof("[%s] Concurrency session created.", consumerID)

	election := concurrency.NewElection(session, constant.Etcd_DatastoreLeaderElectionPath)
	for {
		select {
		case <-a.Ctx.Done():
			logg.L.Infof("[%s] Context cancelled, stopping election campaign.", consumerID)
			return
		default:
			logg.L.Infof("[%s] Attempting to campaign for leadership...", consumerID)
			err := election.Campaign(a.Ctx, consumerID)
			if err != nil {
				logg.L.Infof("[%s] Campaign error: %v. Retrying after delay...", consumerID, err)
				time.Sleep(2 * time.Second)
				continue
			}

			logg.L.Infof("[%s] Successfully elected as leader!", consumerID)
			resp, err := election.Leader(a.Ctx)
			if err != nil {
				logg.L.Errorf("[%s] Failed to get leader key info after elected: %v", consumerID, err)
			} else if resp != nil && len(resp.Kvs) > 0 {
				logg.L.Infof("[%s] Current leader key: %s, revision: %d", consumerID, string(resp.Kvs[0].Key), resp.Kvs[0].ModRevision)
			}

			leaderCtx, leaderCancel := context.WithCancel(a.Ctx)
			leaderLostCh := make(chan struct{})
			go func() {
				defer close(leaderLostCh)
				// Observe 会在当前 Leader key 发生变化时收到通知
				observeCh := election.Observe(leaderCtx)
				select {
				case <-leaderCtx.Done():
					return
				case resp, ok := <-observeCh:
					if ok {
						if string(resp.Kvs[0].Value) != consumerID {
							logg.L.Infof("[%s] Observed leadership change. New leader value: %s. Current leader resigning.", consumerID, string(resp.Kvs[0].Value))
						} else {
							logg.L.Infof("[%s] Observed own leadership key update. Ignoring.", consumerID)
							return
						}
					} else {
						logg.L.Infof("[%s] Leader observation channel closed.", consumerID)
					}
					leaderCancel()
				}
			}()

			a.runLeaderTasks(leaderCtx, consumerID) // 这个函数会阻塞直到 Leader 任务完成或被取消

			leaderCancel() // 确保 leaderCtx 被取消
			<-leaderLostCh

			logg.L.Infof("[%s] Resigning leadership...", consumerID)
			resignCtx, resignCancel := context.WithTimeout(context.Background(), 5*time.Second) // 给 Resign 设置超时
			err = election.Resign(resignCtx)
			resignCancel()
			if err != nil {
				logg.L.Infof("[%s] Failed to resign leadership: %v", consumerID, err)
			} else {
				logg.L.Infof("[%s] Successfully resigned leadership.", consumerID)
			}
			logg.L.Infof("[%s] Re-entering election loop as follower.", consumerID)
		}
	}
}
func (a *app) performReassignment(consumerID string) {
	logg.L.Infof("[%s] Leader performing reassignment...", consumerID)

	vPartitions := defaultVPartitions // Use default first
	// read V from etcd config
	// getVResp, err := client.Get(ctx, virtualPartitionsConf)
	// if err == nil && len(getVResp.Kvs) > 0 {
	//     // Parse value, handle errors
	//     // vPartitions = parsedValue
	// } else if err != nil {
	//    log.Printf("[%s] Warning: Failed to get V partitions from %s: %v. Using default %d", consumerID, virtualPartitionsConf, err, defaultVPartitions)
	// }
	// if vPartitions%4 != 0 {
	//	 logger.Fatalf("[%s] ERROR: Virtual partition count V (%d) must be a multiple of 4. Aborting reassignment.", consumerID, vPartitions)
	// }

	logg.L.Infof("[%s] Using V = %d virtual partitions.", consumerID, vPartitions)

	// 2. Get Active Consumers
	getConsumersResp, err := a.cli.Get(a.Ctx, constant.Etcd_DatastorePrefixConsumerReg, clientv3.WithPrefix())
	if err != nil {
		logg.L.Errorf("[%s] Failed to get active consumers from %s: %v", consumerID, constant.Etcd_DatastorePrefixConsumerReg, err)
		return
	}

	activeConsumers := make([]string, 0, len(getConsumersResp.Kvs))
	for _, kv := range getConsumersResp.Kvs {
		consumerKey := string(kv.Key)
		cID := strings.TrimPrefix(consumerKey, constant.Etcd_DatastorePrefixConsumerReg)
		if cID != "" {
			activeConsumers = append(activeConsumers, cID)
		}
	}

	if len(activeConsumers) == 0 {
		logg.L.Warnf("[%s] No active consumers found. Skipping assignment.", consumerID)
		// deleteExistingAssignments(ctx, client, assignmentPrefix, map[string]struct{}{})
		return
	}

	sort.Strings(activeConsumers)
	logg.L.Infof("[%s] Found %d active consumers: %v", consumerID, len(activeConsumers), activeConsumers)

	// Assignment
	numConsumers := len(activeConsumers)
	partitionsPerSet := vPartitions / 4

	// Divide Consumers into 4 Groups
	consumerGroups := make([][]string, 4)
	baseSize := numConsumers / 4
	remainder := numConsumers % 4
	currentConsumerIndex := 0
	for g := 0; g < 4; g++ {
		groupSize := baseSize
		if g < remainder {
			groupSize++
		}
		if currentConsumerIndex+groupSize > numConsumers { // Boundary check
			groupSize = numConsumers - currentConsumerIndex
		}
		if groupSize > 0 {
			consumerGroups[g] = activeConsumers[currentConsumerIndex : currentConsumerIndex+groupSize]
			currentConsumerIndex += groupSize
		} else {
			consumerGroups[g] = []string{} // Empty group
		}
	}

	newAssignments := make(map[string][]int) // consumerID -> []partitionIDs

	for g := 0; g < 4; g++ { // Iterate through groups/sets
		groupConsumers := consumerGroups[g]
		numConsumersInGroup := len(groupConsumers)

		if numConsumersInGroup == 0 {
			continue // Skip empty groups
		}

		startPartition := g * partitionsPerSet
		endPartition := (g + 1) * partitionsPerSet // Exclusive end for loop

		partitionIndexInSet := 0
		for p := startPartition; p < endPartition; p++ {
			targetConsumerIndex := partitionIndexInSet % numConsumersInGroup
			targetConsumerID := groupConsumers[targetConsumerIndex]

			if _, ok := newAssignments[targetConsumerID]; !ok {
				newAssignments[targetConsumerID] = make([]int, 0, partitionsPerSet) // Estimate capacity
			}
			newAssignments[targetConsumerID] = append(newAssignments[targetConsumerID], p)
			partitionIndexInSet++
		}
	}

	logg.L.Infof("[%s] Writing %d assignments to etcd...", consumerID, len(newAssignments))
	// Use a transaction or batch puts for better efficiency/atomicity if needed
	var ops []clientv3.Op
	assignedKeys := make(map[string]struct{}) // Keep track of keys we just assigned

	for cID, partitions := range newAssignments {
		assignmentKey := constant.Etcd_DatastoreAssignmentPrefix + cID
		assignedKeys[assignmentKey] = struct{}{} // Mark this key as actively assigned

		// Sort partitions for consistent representation in etcd (optional but good)
		sort.Ints(partitions)

		jsonData, err := sonic.Marshal(partitions)
		if err != nil {
			logg.L.Errorf("[%s] Failed to marshal assignment list for %s: %v", consumerID, cID, err)
			continue
		}

		ops = append(ops, clientv3.OpPut(assignmentKey, string(jsonData)))
	}

	getExistingAssignResp, err := a.cli.Get(a.Ctx, constant.Etcd_DatastoreAssignmentPrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		logg.L.Infof("[%s] Warning: Failed to get existing assignments for cleanup: %v", consumerID, err)
	} else {
		for _, kv := range getExistingAssignResp.Kvs {
			existingKey := string(kv.Key)
			if _, isActive := assignedKeys[existingKey]; !isActive {
				// This key exists in etcd but wasn't part of the new assignment
				logg.L.Infof("[%s] Deleting stale assignment key: %s", consumerID, existingKey)
				ops = append(ops, clientv3.OpDelete(existingKey))
			}
		}
	}

	// Execute all Puts (and Deletes) in a single transaction (more atomic)
	if len(ops) > 0 {
		txn := a.cli.Txn(a.Ctx)
		txnResp, err := txn.Then(ops...).Commit()
		if err != nil {
			logg.L.Errorf("[%s] Failed to commit assignment transaction: %v", consumerID, err)
		} else if !txnResp.Succeeded {
			logg.L.Errorf("[%s] Assignment transaction failed (maybe contention?)", consumerID)
		} else {
			logg.L.Infof("[%s] Successfully committed %d assignment operations.", consumerID, len(ops))
		}
	} else {
		logg.L.Warnf("[%s] No assignment operations needed.", consumerID)
	}

	logg.L.Infof("[%s] Reassignment finished.", consumerID)
}
