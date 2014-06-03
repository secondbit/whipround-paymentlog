package paymentlog

import (
	"sync"
	"time"
)

type MemoryStore struct {
	paymentLogs map[string]*PaymentLog
	failureLogs map[string]*FailureLog
	sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		paymentLogs: make(map[string]*PaymentLog),
		failureLogs: make(map[string]*FailureLog),
	}
}

func (store *MemoryStore) StorePaymentLog(log PaymentLog) error {
	store.Lock()
	defer store.Unlock()
	if _, ok := store.paymentLogs[log.ID]; ok {
		return AlreadyExists
	}
	store.paymentLogs[log.ID] = &log
	return nil
}

func (store *MemoryStore) UpdatePaymentLog(id string, change PaymentLogChange) error {
	store.Lock()
	defer store.Unlock()
	if _, ok := store.paymentLogs[id]; !ok {
		return LogNotFound
	}
	if change.Amount != nil {
		store.paymentLogs[id].Amount = *change.Amount
	}
	if change.Description != nil {
		store.paymentLogs[id].Description = *change.Description
	}
	if change.Source != nil {
		store.paymentLogs[id].Source = *change.Source
	}
	if change.SourceID != nil {
		store.paymentLogs[id].SourceID = *change.SourceID
	}
	if change.Created != nil {
		store.paymentLogs[id].Created = *change.Created
	}
	if change.Updated != nil {
		store.paymentLogs[id].Updated = *change.Updated
	}
	if change.Status != nil {
		store.paymentLogs[id].Status = *change.Status
	}
	if change.Currency != nil {
		store.paymentLogs[id].Currency = *change.Currency
	}
	return nil
}

func (store *MemoryStore) DeletePaymentLog(id string) error {
	store.Lock()
	defer store.Unlock()
	if _, ok := store.paymentLogs[id]; !ok {
		return LogNotFound
	}
	delete(store.paymentLogs, id)
	return nil
}

func (store *MemoryStore) GetPaymentLog(id string) (PaymentLog, error) {
	store.Lock()
	defer store.Unlock()
	if log, ok := store.paymentLogs[id]; !ok || log == nil {
		return PaymentLog{}, LogNotFound
	} else {
		return *log, nil
	}
}

func (store *MemoryStore) ListPaymentLogsByProject(id string, num, offset int) ([]PaymentLog, error) {
	store.Lock()
	defer store.Unlock()
	results := make([]PaymentLog, 0)
	for _, log := range store.paymentLogs {
		if log == nil {
			continue
		}
		if log.ProjectID == id {
			results = append(results, *log)
		}
	}
	return SortLogsByCreated(results), nil
}

func (store *MemoryStore) ListPaymentLogsByUser(id string, num, offset int) ([]PaymentLog, error) {
	store.Lock()
	defer store.Unlock()
	results := make([]PaymentLog, 0)
	for _, log := range store.paymentLogs {
		if log == nil {
			continue
		}
		if log.UserID == id {
			results = append(results, *log)
		}
	}
	return SortLogsByCreated(results), nil
}

func (store *MemoryStore) ListPaymentLogs(num, offset int) ([]PaymentLog, error) {
	store.Lock()
	defer store.Unlock()
	results := make([]PaymentLog, 0)
	for _, log := range store.paymentLogs {
		if log == nil {
			continue
		}
		results = append(results, *log)
	}
	return SortLogsByCreated(results), nil
}

func (store *MemoryStore) StoreFailureLog(log FailureLog) error {
	store.Lock()
	defer store.Unlock()
	if _, ok := store.failureLogs[log.ID]; ok {
		return AlreadyExists
	}
	store.failureLogs[log.ID] = &log
	return nil
}

func (store *MemoryStore) ListFailureLogs(num, offset int) ([]FailureLog, error) {
	store.Lock()
	defer store.Unlock()
	results := make([]FailureLog, 0)
	for _, log := range store.failureLogs {
		if log == nil {
			continue
		}
		results = append(results, *log)
	}
	return SortFailureLogs(results), nil
}

func (store *MemoryStore) ListFailureLogsSince(timestamp time.Time) ([]FailureLog, error) {
	store.Lock()
	defer store.Unlock()
	results := make([]FailureLog, 0)
	for _, log := range store.failureLogs {
		if log == nil {
			continue
		}
		if !log.Timestamp.After(timestamp) {
			continue
		}
		results = append(results, *log)
	}
	return SortFailureLogs(results), nil
}
