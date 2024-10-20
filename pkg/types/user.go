package types

type Role uint8

const (
	RoleVisitor Role = iota
	RoleUser
	RoleCore
	RoleGateway
	RoleAgent
)
