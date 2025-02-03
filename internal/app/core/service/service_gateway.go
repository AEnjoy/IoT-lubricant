package service

import "github.com/AEnjoy/IoT-lubricant/internal/model/repo"

type GatewayService struct {
	db repo.CoreDbOperator
}
type IGatewayService interface {
}
