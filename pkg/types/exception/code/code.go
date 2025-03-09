package code

type ResCode int64

// Common
const (
	Success ResCode = 100000 + iota
	_
	ErrorBadRequest
	ErrorNotFound
	ErrorInternalServerError
	ErrorUnknown
	ErrorIO
	DeadLine
	ErrorPushTaskFailed
)

// Auth
const (
	ErrorTokenInvalid ResCode = 100101 + iota
	ErrorTokenExpired
	ErrorInvalidAuthHeader
	ErrorInvalidAuthKey
	ErrorForbidden
	ErrorGetClaimsFailed
)

const (
	ErrorBind ResCode = 100201 + iota
	ErrorEncodeFailed
	ErrorDecodeFailed
	ErrorEncodeJSON
	ErrorDecodeJSON
	ErrorValidation
	ErrorReadRequestBody
	ErrorEncodeProtoMessage
	ErrorDecodeProtoMessage
)

// core module
const (
	ErrorCoreNoTask ResCode = 120001 + iota
	ErrorCoreTaskTimeout
	DbAddGatewayFailed
	DbUpdateGatewayInfoFailed
	DbGetGatewayFailed
	LinkToGatewayFailed
	ErrorDeployGatewayFailed
	GetGatewayFailed
	AddGatewayFailed
	AddGatewayHostFailed
	RemoveGatewayFailed
	RemoveGatewayHostFailed
	DescriptionHostFailed
	ErrorCommunicationWithAuthServer
	ErrorGetGatewayStatusFailed
)

// gateway
const (
	ErrorGatewayAgentNotFound ResCode = 130001 + iota
	WarnAgentOffline
	ErrGaterDataReqFailed
	OperationOnlyAtLocal
	ErrorAgentStartFailed
	AddAgentFailed
)

// agent
const (
	ErrorAgentInvalidConfig ResCode = 140001 + iota
	ErrorAgentNotAllowMultiGatherInstance
	ErrorAgentNeedInit
	GetAgentFailed
	UpdateAgentFailed
	SetAgentFailed
	StartAgentFailed
	StopAgentFailed
	RemoveAgentFailed
	ErrorNoAgentContainerConfSet
	ErrorAgentUpdateFailed
	ErrorAgentUpdateNotSupportRemote
	ErrGaterStartFailed
)

// cache
const (
	ErrorCacheNeedInit ResCode = 200001 + iota
	ErrorCacheNullCache
)

// database
const (
	ErrorDbNeedTxn ResCode = 210001 + iota
)

// openapi
const (
	ErrorApiNotFound ResCode = 220001 + iota
	ErrorApiInvalidMethod
	ErrorApiInvalidInput
	ErrorApiInvalidPath
	ErrorApiInvalidSlot
	ErrorApiNotInit
)

// docker
const (
	ErrContainerNotRunning ResCode = 230001 + iota
)

// mq
const (
	MqPublishFailed ResCode = 240001 + iota
	MqSubscribeFailed
)

// request parameters
const (
	ErrorGatewayHostNeedPasswdOrPrivateKey ResCode = 410001 + iota
)

var StatusMsgMap = map[ResCode]string{
	// Common
	Success:                  "success",
	ErrorBadRequest:          "Invalid Request",
	ErrorNotFound:            "Not Found",
	ErrorInternalServerError: "Internal Service Error",
	ErrorUnknown:             "Unknown Error",
	ErrorTokenInvalid:        "Invalid Token",
	ErrorTokenExpired:        "Token Expired",
	ErrorBind:                "Error occurred when binding the request body",
	ErrorReadRequestBody:     "Error occurred when reading the request body",
	ErrorEncodeFailed:        "Encoding failed due to an error with the data",
	ErrorDecodeFailed:        "Decoding failed due to an error with the data",
	ErrorEncodeJSON:          "JSON encode failed",
	ErrorDecodeJSON:          "JSON decode failed",
	ErrorEncodeProtoMessage:  "Protobuf encode failed",
	ErrorDecodeProtoMessage:  "Protobuf decode failed",
	ErrorValidation:          "Validation Error",
	ErrorInvalidAuthHeader:   "Invalid authorization header",
	ErrorInvalidAuthKey:      "Invalid authorization key",
	ErrorForbidden:           "Permission Denied",
	ErrorIO:                  "IO error",
	ErrorGetClaimsFailed:     "Get claims(user information) from context failed",
	DeadLine:                 "context deadline or cancel",
	ErrorPushTaskFailed:      "push task failed",

	// Core
	ErrorCoreNoTask:                  "target has no task",
	ErrorCoreTaskTimeout:             "get task timeout",
	AddGatewayHostFailed:             "add gateway host failed",
	AddGatewayFailed:                 "add gateway failed",
	RemoveGatewayFailed:              "remove gateway failed",
	RemoveGatewayHostFailed:          "remove gateway host failed",
	DescriptionHostFailed:            "get description host failed",
	ErrorCommunicationWithAuthServer: "communication with auth service failed",
	ErrorGetGatewayStatusFailed:      "get gateway status failed",

	// Gateway
	ErrorGatewayAgentNotFound: "agent not found",
	WarnAgentOffline:          "agent is offline",
	ErrGaterDataReqFailed:     "gather data request failed",
	OperationOnlyAtLocal:      "only supports local agents",
	ErrorAgentStartFailed:     "agent start failed",
	AddAgentFailed:            "add agent failed",

	// Agent
	ErrorAgentInvalidConfig:               "invalid config",
	ErrorAgentNotAllowMultiGatherInstance: "not allow multi gather instance",
	ErrorAgentNeedInit:                    "should be call lubricant.agent.edgeService / setAgent before this operation",
	GetAgentFailed:                        "get agent failed",
	UpdateAgentFailed:                     "update agent failed",
	SetAgentFailed:                        "set agent failed",
	StartAgentFailed:                      "start agent failed",
	StopAgentFailed:                       "stop agent failed",
	RemoveAgentFailed:                     "remove agent failed",
	ErrorNoAgentContainerConfSet:          "agent container conf is not set",
	ErrorAgentUpdateFailed:                "update agent failed",
	ErrorAgentUpdateNotSupportRemote:      "update agent operation not support remote agent",
	ErrGaterStartFailed:                   "agent start gather failed",

	// Cache
	ErrorCacheNeedInit:  "cache client need init",
	ErrorCacheNullCache: "cache client is nil",

	// Database
	ErrorDbNeedTxn: "this operation need start with txn support",

	// Openapi
	ErrorApiNotFound:      "not found",
	ErrorApiInvalidMethod: "invalid method",
	ErrorApiInvalidInput:  "invalid input",
	ErrorApiInvalidPath:   "invalid path",
	ErrorApiInvalidSlot:   "invalid slot",
	ErrorApiNotInit:       "not initialized",

	// Docker
	ErrContainerNotRunning: "container is not running",

	// MQ
	MqPublishFailed:   "publish message to messageQueue failed",
	MqSubscribeFailed: "subscribe message from messageQueue failed",

	// Request parameters
	ErrorGatewayHostNeedPasswdOrPrivateKey: "gateway host need passwd or private key for remote login(ssh)",
}

// GetMsg 返回状态码对应msg
func (c ResCode) GetMsg() string {
	msg, ok := StatusMsgMap[c]
	if !ok {
		return StatusMsgMap[ErrorUnknown]
	}
	return msg
}
