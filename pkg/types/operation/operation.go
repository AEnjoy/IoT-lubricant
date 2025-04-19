package operation

import "encoding/json"

type Operation uint8

const (
	_ Operation = 0
	OperationNil
	OperationUserLogin             Operation = 10
	OperationUserLogout            Operation = 11
	OperationUserChangePassword    Operation = 12
	OperationAddTask               Operation = 13
	OperationQueryTask             Operation = 14
	OperationViewTaskResult        Operation = 15
	OperationAddGateway            Operation = 20
	OperationRemoveGateway         Operation = 21
	OperationAddAgent              Operation = 22
	OperationRemoveAgent           Operation = 23
	OperationAddSchedule           Operation = 24
	OperationRemoveSchedule        Operation = 25
	OperationAddDriverContainer    Operation = 30
	OperationRemoveDriverContainer Operation = 31
	OperationAddAgentContainer     Operation = 32
	OperationRemoveAgentContainer  Operation = 33
	OperationEnableOpenAPI         Operation = 40
	OperationDisableOpenAPI        Operation = 41
	OperationSendRequest           Operation = 42
	OperationGetOpenAPIDoc         Operation = 43
	OperationGetEnableOpenAPI      Operation = 44
)

type Operator interface {
	TaskOperation() Operation
	json.Marshaler
}
