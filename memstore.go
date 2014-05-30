package paymentlog

import (
	"sync"
)

type MemoryStore struct {
	paymentLogs map[string]PaymentLog
	failureLogs map[string]FailureLog
	sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		paymentLogs: make(map[string]PaymentLog),
		failureLogs: make(map[string]FailureLog),
	}
}

func (store *MemoryStore) StorePaymentLog(log PaymentLog) error {
	store.Lock()
	defer store.Unlock()
	if _, ok := store.paymentLogs[log.ID]; ok {
		return AlreadyExists
	}
	store.paymentLogs[log.ID] = log
	return nil
}
