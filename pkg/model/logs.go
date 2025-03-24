package model

type Log struct {
	ID         int    `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	LogID      string `json:"logId" gorm:"column:log_id"`
	OperatorID string `json:"operatorId" gorm:"column:operator_id"` // UserID / DevicesID
	// todo
}
