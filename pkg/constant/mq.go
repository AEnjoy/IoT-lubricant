package constant

// logger
const (
	MESSAGE_SVC_LOGGER      = "/logger/message"
	DATASTORE_PROJECT       = "/handler/data/projectid" // todo: refactor: 应该使用 projectId 而不是 userId
	DATASTORE_USER_REG_RESP = DATASTORE_PROJECT + "/%s/reg/response"
	DATASTORE_USER_CLOSE    = DATASTORE_PROJECT + "/%s/close"
	DATASTORE_PROJECT_DATA  = DATASTORE_PROJECT + "/%s/data"
)
