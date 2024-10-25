package types

import "time"

type Device struct {
	Id     string `json:"id" gorm:"column:id;primary_key"`
	UserId string `json:"user_id" gorm:"column:user_id"`

	DeviceBasicInfo

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Device) TableName() string {
	return "devices"
}

type DeviceBasicInfo struct {
	Name            string `json:"name" gorm:"column:name" validate:"required,max=64"`
	Type            string `json:"type" gorm:"column:type" validate:"required,max=64"`
	OperationSystem string `json:"os" gorm:"column:os" validate:"required,max=64"`
	Manufacturer    string `json:"manufacturer" gorm:"column:manufacturer" validate:"required,max=64"`
	Model           string `json:"model" gorm:"column:model" validate:"required,max=64"`
	Protocol        string `json:"protocol" gorm:"column:protocol" validate:"required,max=64"`
	Language        string `json:"language" gorm:"column:language" validate:"required,max=64"`
}

func (DeviceBasicInfo) TableName() string {
	return "device_basic_info"
}

type DeviceAPI struct {
	Method      string `json:"method" validate:"required,oneof=GET POST DELETE PUT"`
	Path        string `json:"path" validate:"path"`
	Description string `json:"description" validate:"required,max=1024"`
}

func (DeviceAPI) TableName() string {
	return "device_api"
}
