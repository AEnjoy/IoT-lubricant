package constant

// logger
const (
	MESSAGE_SVC_LOGGER      = "/logger/message"
	DATASTORE_PROJECT       = "/handler/data/projectid"
	DATASTORE_USER_REG_RESP = DATASTORE_PROJECT + "/%s/reg/response"
	DATASTORE_USER_CLOSE    = DATASTORE_PROJECT + "/%s/close"
	DATASTORE_PROJECT_DATA  = DATASTORE_PROJECT + "/%s/data"
)
