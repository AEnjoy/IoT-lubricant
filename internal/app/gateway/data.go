package gateway

import (
	"context"
	"errors"

	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/google/uuid"
)

//	func (a *app) agentHandelSignal(id string) {
//		//send
//		go func() {
//			if err := a.agentPushDataToServer(a.ctrl, id); err != nil {
//				logger.Error(err)
//			}
//		}()
//		go func() {
//			if err := a.handelSendSignal(id); err != nil {
//				logger.Error(err)
//			}
//		}()
//
//		//todo:recv
//	}
//
//	func (a *app) handelSendSignal(id string) error {
//		// todo: not all implemented yet
//		//  task: 1. need to support choose report cycle
//		//        2. need to support modify and trigger by core server manually
//		v, ok := agentDataStore.Load(id)
//		if !ok {
//			return errs.ErrAgentNotFound
//		}
//
//		agentMap := v.(*agentData)
//		reportCycle := agentPool[id].agentInfo.Cycle
//		for range time.Tick(time.Second * time.Duration(reportCycle)) {
//			agentMap.sendSignal <- struct{}{}
//		}
//		return nil
//	}
//
//	func (a *app) agentPushDataToServer(ctx context.Context, id string) error {
//		v, ok := agentDataStore.Load(id)
//		if !ok {
//			return errs.ErrAgentNotFound
//		}
//		agentMap := v.(*agentData)
//		for {
//			select {
//			case <-ctx.Done():
//				return nil
//			case <-agentMap.sendSignal:
//				go func() {
//					_ = pushDataToServer(ctx, agentMap, a.grpcClient, id)
//					// todo: this error need to be handled
//				}()
//			}
//		}
//	}
func pushDataToServer(ctx context.Context, agentMap *agentData, grpcClient core.CoreServiceClient, id string) error {
	stream, err := grpcClient.PushDataStream(ctx)
	if err != nil {
		return err
	}
	data := agentMap.coverToGrpcData()
	if data == nil {
		return nil
	}
	agentMap.cleanData()

	data.GatewayId = gatewayId
	data.AgentID = id
	data.MessageId = uuid.NewString()

	err = stream.Send(data)
	if err != nil {
		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	if resp.GetMessageId() != data.MessageId {
		return errors.New("message id not match")
	}
	return nil
}
