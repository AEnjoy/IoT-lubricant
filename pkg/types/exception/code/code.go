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
)

// Auth
const (
	ErrorTokenInvalid ResCode = 100101 + iota
	ErrorTokenExpired
	ErrorInvalidAuthHeader
	ErrorInvalidAuthKey
	ErrorForbidden
)

const (
	ErrorBind ResCode = 100201 + iota
	ErrorEncodeFailed
	ErrorDecodeFailed
	ErrorEncodeJSON
	ErrorDecodeJSON
	ErrorValidation
	ErrorReadRequestBody
)

// core module
const (
	ErrorCoreNoTask ResCode = 120001 + iota
	ErrorCoreTaskTimeout
)

// gateway
const (
	ErrorGatewayAgentNotFound ResCode = 130001 + iota
)

// agent
const (
	ErrorAgentInvalidConfig ResCode = 140001 + iota
	ErrorAgentNotAllowMultiGatherInstance
	ErrorAgentNeedInit
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

var StatusMsgMap = map[ResCode]string{
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
	ErrorValidation:          "Validation Error",
	ErrorInvalidAuthHeader:   "Invalid authorization header",
	ErrorInvalidAuthKey:      "Invalid authorization key",
	ErrorForbidden:           "Permission Denied",

	ErrorCoreNoTask:      "target has no task",
	ErrorCoreTaskTimeout: "get task timeout",

	ErrorGatewayAgentNotFound: "agent not found",

	ErrorAgentInvalidConfig:               "invalid config",
	ErrorAgentNotAllowMultiGatherInstance: "not allow multi gather instance",
	ErrorAgentNeedInit:                    "should be call lubricant.agent.edgeService / setAgent before this operation",

	ErrorCacheNeedInit:  "cache client need init",
	ErrorCacheNullCache: "cache client is nil",

	ErrorDbNeedTxn: "this operation need start with txn support",

	ErrorApiNotFound:      "not found",
	ErrorApiInvalidMethod: "invalid method",
	ErrorApiInvalidInput:  "invalid input",
	ErrorApiInvalidPath:   "invalid path",
	ErrorApiInvalidSlot:   "invalid slot",
	ErrorApiNotInit:       "not initialized",
}

// GetMsg 返回状态码对应msg
func (c ResCode) GetMsg() string {
	msg, ok := StatusMsgMap[c]
	if !ok {
		return StatusMsgMap[ErrorUnknown]
	}
	return msg
}
