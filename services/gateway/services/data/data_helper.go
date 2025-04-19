package data

// ManualPushAgentData 该操作用于销毁agentApiControl前最后一次的手动向dataCollect发送数据
func ManualPushAgentData(id string) {
	if v, ok := agentDataStore.LoadAndDelete(id); ok {
		api := v.(Apis)
		data := api.Pop()
		if data != nil {
			data.AgentID = id
			dataSendQueue <- data
		}
	}
}
