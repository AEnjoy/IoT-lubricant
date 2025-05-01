package constant

// logger
const (
	MESSAGE_SVC_LOGGER      = "/logger/message"
	DATASTORE_USER          = "/handler/data/userid" // todo: refactor: 应该使用 projectId 而不是 userId
	DATASTORE_USER_REG_RESP = DATASTORE_USER + "/%s/reg/response"
	DATASTORE_USER_CLOSE    = DATASTORE_USER + "/%s/close"
	DATASTORE_USER_DATA     = DATASTORE_USER + "/%s/data"
)
