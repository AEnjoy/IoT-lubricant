package mq

import (
	"context"
)

// GatewayOfflineSignal: 监听网关下线信号, 如果 err != nil 则表示网关是异常下线
//
// Deprecated :
func (m *MqService) GatewayOfflineSignal(ctx context.Context, gatewayID string, err error) error {
	panic("this function is deprecated need clean")
	//defer func() {
	//	txn := m.DataStore.ICoreDb.Begin()
	//	_ = m.DataStore.ICoreDb.SetGatewayStatus(ctx, txn, gatewayID, "offline")
	//	m.DataStore.ICoreDb.Commit(txn)
	//}()
	//
	//if err != nil {
	//	return m.Mq.Publish(fmt.Sprintf("/monitor/%s/%s/offline/error", taskTypes.TargetGateway, gatewayID),
	//		[]byte(fmt.Sprintf("Time:%s,Err:%v", time.Now().Format("2006-01-02 15:04:05"), err)))
	//}
	//return m.Mq.Publish(fmt.Sprintf("/monitor/%s/%s/offline", taskTypes.TargetGateway, gatewayID),
	//	[]byte(time.Now().Format("2006-01-02 15:04:05")))
}
