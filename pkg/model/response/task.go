package response

type QueryTaskResultResponse struct {
	TaskID string `json:"taskId"`
	Status string `json:"status"`
	Result string `json:"result"`
}
