package services

import (
	errCh "github.com/aenjoy/iot-lubricant/pkg/error"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
	"gorm.io/gorm"
)

// txnHelper return: database txn, error channel,and a function that must call by defer
func (s *GatewayService) txnHelper() (txn *gorm.DB, errorCh *errCh.ErrorChan, f func()) {
	return _txnHelper(s.db)
}

// txnHelper return: database txn, error channel,and a function that must call by defer
func (s *AgentService) txnHelper() (txn *gorm.DB, errorCh *errCh.ErrorChan, f func()) {
	return _txnHelper(s.db)
}
func _txnHelper(db repo.ICoreDb) (txn *gorm.DB, errorCh *errCh.ErrorChan, f func()) {
	txn = db.Begin()
	errorCh = errCh.NewErrorChan()
	f = errCh.HandleErrorCh(errorCh).
		ErrorWillDo(func(error) {
			db.Rollback(txn)
		}).
		SuccessWillDo(func() {
			db.Commit(txn)
		}).
		Do
	return
}
