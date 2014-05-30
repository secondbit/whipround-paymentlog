package paymentlog

import (
	"testing"
	"time"
)

func comparePaymentLogs(expectation, result PaymentLog, t *testing.T) (success bool, field string, expectedValue, resultValue interface{}) {
	if expectation.ID != result.ID {
		return false, "id", expectation.ID, result.ID
	}
	if expectation.Amount != result.Amount {
		return false, "amount", expectation.Amount, result.Amount
	}
	if expectation.Description != result.Description {
		return false, "description", expectation.Description, result.Description
	}
	if expectation.Source != result.Source {
		return false, "source", expectation.Source, result.Source
	}
	if expectation.SourceID != result.SourceID {
		return false, "source ID", expectation.SourceID, result.SourceID
	}
	if !expectation.Created.Equal(result.Created) {
		return false, "created", expectation.Created, result.Created
	}
	if !expectation.Updated.Equal(result.Updated) {
		return false, "updated", expectation.Updated, result.Updated
	}
	if expectation.Status != result.Status {
		return false, "status", expectation.Status, result.Status
	}
	if expectation.Currency != result.Currency {
		return false, "currency", expectation.Currency, result.Currency
	}
	if expectation.CampaignID != result.CampaignID {
		return false, "campaign ID", expectation.CampaignID, result.CampaignID
	}
	if expectation.GoalID != result.GoalID {
		return false, "goal ID", expectation.GoalID, result.GoalID
	}
	if expectation.UserID != result.UserID {
		return false, "user ID", expectation.UserID, result.UserID
	}
	if expectation.AccountID != result.AccountID {
		return false, "account ID", expectation.AccountID, result.AccountID
	}
	if expectation.AccountType != result.AccountType {
		return false, "account type", expectation.AccountType, result.AccountType
	}
	return true, "", nil, nil
}

func TestStoringPaymentLogInMemory(t *testing.T) {
	store := NewMemoryStore()
	p := PaymentLog{
		ID:          "test-payment-log",
		Amount:      1,
		Source:      SourceBalanced,
		SourceID:    "balanced-id",
		Created:     time.Now(),
		Status:      StatusPending,
		Currency:    CurrencyUSD,
		CampaignID:  "campaign-id",
		GoalID:      "goal-id",
		UserID:      "user-id",
		AccountID:   "account-id",
		AccountType: "google",
	}
	err := store.StorePaymentLog(p)
	if err != nil {
		t.Errorf("Error storing payment log in memory: %s", err)
	}
	p2, ok := store.paymentLogs[p.ID]
	if !ok {
		t.Errorf("PaymentLog never got stored in memory: %+v", store.paymentLogs)
	}
	success, field, expectation, result := comparePaymentLogs(p, p2, t)
	if !success {
		t.Errorf("Mismatch. Expected payment log %s to be %+v, got %+v.", field, expectation, result)
	}
}

func TestStoringDuplicatePaymentLogInMemory(t *testing.T) {
	store := NewMemoryStore()
	p := PaymentLog{
		ID:          "test-payment-log",
		Amount:      1,
		Source:      SourceBalanced,
		SourceID:    "balanced-id",
		Created:     time.Now(),
		Status:      StatusPending,
		Currency:    CurrencyUSD,
		CampaignID:  "campaign-id",
		GoalID:      "goal-id",
		UserID:      "user-id",
		AccountID:   "account-id",
		AccountType: "google",
	}
	err := store.StorePaymentLog(p)
	if err != nil {
		t.Errorf("Error storing payment log in memory: %s", err)
	}
	err = store.StorePaymentLog(p)
	if err != AlreadyExists {
		t.Errorf("Expected %s when storing payment log in memory, got %v", AlreadyExists, err)
	}
}
