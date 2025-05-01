package datastoreAssistant

import (
	"context"
	"errors"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const leaseTTL = 10

func (a *app) registerConsumer(consumerID string) (clientv3.LeaseID, error) {
	regKey := constant.Etcd_DatastorePrefixConsumerReg + consumerID
	regValue := consumerID

	leaseResp, err := a.cli.Grant(a.Ctx, leaseTTL)
	if err != nil {
		logg.L.Errorf("[%s] Failed to grant lease: %v", consumerID, err)
		return 0, err
	}
	leaseID := leaseResp.ID
	logg.L.Errorf("[%s] Granted lease ID: %d", consumerID, leaseID)

	_, err = a.cli.Put(a.Ctx, regKey, regValue, clientv3.WithLease(leaseID))
	if err != nil {
		logg.L.Errorf("[%s] Failed to register consumer with lease %d: %v", consumerID, leaseID, err)

		_, err2 := a.cli.Revoke(a.Ctx, leaseID)
		return 0, errors.Join(err, err2)
	}
	logg.L.Infof("[%s] Registered consumer at %s with lease %d", consumerID, regKey, leaseID)

	keepAliveCh, err := a.cli.KeepAlive(a.Ctx, leaseID)
	if err != nil {
		logg.L.Errorf("[%s] Failed to start keepalive for lease %d: %v", consumerID, leaseID, err)
		_, err2 := a.cli.Revoke(a.Ctx, leaseID)
		return 0, errors.Join(err, err2)
	}

	// 处理 KeepAlive request
	go func() {
		for {
			select {
			case _, ok := <-keepAliveCh:
				if !ok {
					logg.L.Warnf("[%s] Keepalive channel closed for lease %d.", consumerID, leaseID)
					return
				}
			case <-a.Ctx.Done():
				logg.L.Infof("[%s] Context cancelled, stopping keepalive for lease %d.", consumerID, leaseID)
				revokeCtx, revokeCancel := context.WithTimeout(context.Background(), 2*time.Second)
				_, err := a.cli.Revoke(revokeCtx, leaseID)
				if err != nil {
					logg.L.Errorf("[%s] Failed to revoke lease %d: %v", consumerID, leaseID, err)
				}
				revokeCancel()
				return
			}
		}
	}()

	return leaseID, nil
}

// watchAssignments watches the specific assignment key for this consumer instance
func (a *app) watchAssignments(consumerID string) {
	assignmentKey := constant.Etcd_DatastoreAssignmentPrefix + consumerID
	logg.L.Infof("[%s] Starting to watch assignments at %s", consumerID, assignmentKey)

	getResp, err := a.cli.Get(a.Ctx, assignmentKey)
	if err != nil {
		logg.L.Errorf("[%s] Failed to get initial assignment state from %s: %v", consumerID, assignmentKey, err)
		// Continue to watch, maybe the key will be created later
	} else {
		if len(getResp.Kvs) > 0 {
			logg.L.Infof("[%s] Got initial assignment.", consumerID)
			updateAssignedPartitions(getResp.Kvs[0].Value, consumerID)
		} else {
			logg.L.Infof("[%s] No initial assignment found.", consumerID)
			clearAssignedPartitions(consumerID) // Ensure state is clear initially
		}
	}

	// Determine revision to start watching from
	watchRevision := int64(0) // Start from current revision by default
	if getResp != nil {
		watchRevision = getResp.Header.Revision + 1 // Start watch after the Get operation
	}

	// 2. Start Watch
	watchChan := a.cli.Watch(a.Ctx, assignmentKey, clientv3.WithRev(watchRevision))

	for {
		select {
		case <-a.Ctx.Done():
			logg.L.Warnf("[%s] Stopping assignment watch due to context cancellation.", consumerID)
			return
		case watchResp, ok := <-watchChan:
			if !ok {
				logg.L.Warnf("[%s] Assignment watch channel closed.", consumerID)
				return
			}
			if err := watchResp.Err(); err != nil {
				logg.L.Warnf("[%s] Assignment watch error: %v", consumerID, err)
				time.Sleep(2 * time.Second)       // Avoid tight loop on persistent errors
				go a.watchAssignments(consumerID) // retry
				return
			}

			for _, event := range watchResp.Events {
				logg.L.Infof("[%s] Assignment watch event: Type[%s] Key[%s]", consumerID, event.Type, string(event.Kv.Key))
				switch event.Type {
				case clientv3.EventTypePut: // Assignment created or updated
					updateAssignedPartitions(event.Kv.Value, consumerID)
				case clientv3.EventTypeDelete: // Assignment removed
					clearAssignedPartitions(consumerID)
				}
			}
		}
	}
}
