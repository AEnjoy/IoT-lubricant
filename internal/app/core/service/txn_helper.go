package service

import (
	errCh "github.com/AEnjoy/IoT-lubricant/pkg/error"
	"gorm.io/gorm"
)

// todo: need to merge those two txn function

// txnHelper return: database txn, error channel,and a function that must call by defer
func (s *GatewayService) txnHelper() (txn *gorm.DB, errorCh *errCh.ErrorChan, f func()) {
	txn = s.db.Begin()
	errorCh = errCh.NewErrorChan()
	f = errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		s.db.Rollback(txn)
	}).SuccessWillDo(func() {
		s.db.Commit(txn)
	}).Do
	return
}
func (s *AgentService) txnHelper() (txn *gorm.DB, errorCh *errCh.ErrorChan, f func()) {
	txn = s.db.Begin()
	errorCh = errCh.NewErrorChan()
	f = errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		s.db.Rollback(txn)
	}).SuccessWillDo(func() {
		s.db.Commit(txn)
	}).Do
	return
}
