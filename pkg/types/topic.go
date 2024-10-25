package types

// messageQueue topic list

// Common
const (
	Topic_Ping = "/ping"
)

// Gateway
const (
	Topic_AgentRegister    = "/agent/register/"           // + agentId agent->gateway
	Topic_AgentRegisterAck = "/agent/register/ack/"       // + agentId  gateway->agent
	Topic_AgentDevice      = "/agent/"                    // + agentId agent<->gateway
	Topic_GatewayInfo      = "/gateway/info"              // agent->gateway->agent
	Topic_GatewayData      = "/gateway/data"              // agent->gateway->agent
	Topic_AgentDataPush    = "/gateway/data/push/"        // + agentId agent->gateway
	Topic_AgentDataPushAck = "/gateway/data/push/ack/"    // + agentId agent->gateway
	Topic_MessagePush      = "/gateway/message/push/"     // agent->gateway
	Topic_MessagePushAck   = "/gateway/message/push/ack/" // + messageId gateway->agent
	Topic_MessagePull      = "/gateway/message/pull/"     // + messageId gateway->agent
)

type Ping struct {
	Status  int    `json:"status"`
	Message string `json:"message"` // optional
}
type Register struct {
	ID string `json:"id"`
}
type Command struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

const (
	Command_nil = 0
	Command_Add = iota

	Command_RemoveAgent = 18
)
