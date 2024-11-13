package code

type ResCode int64

// Common
const (
	Success ResCode = 100000 + iota
	_
	ErrorBadRequest
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

var StatusMsgMap = map[ResCode]string{
	Success:                  "success",
	ErrorBadRequest:          "Invalid Request",
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
}

// GetMsg 返回状态码对应msg
func (c ResCode) GetMsg() string {
	msg, ok := StatusMsgMap[c]
	if !ok {
		return StatusMsgMap[ErrorUnknown]
	}
	return msg
}
