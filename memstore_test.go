package paymentlog

import (
	"testing"
	"time"
)

func comparePaymentLogs(expectation, result PaymentLog) (success bool, field string, expectedValue, resultValue interface{}) {
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
	success, field, expectation, result := comparePaymentLogs(p, *p2)
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

func TestUpdatingPaymentLogInMemory(t *testing.T) {
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
	store.paymentLogs[p.ID] = &p
	p.Amount = 2
	p.Description = "new description"
	p.Source = "new source"
	p.SourceID = "new source id"
	p.Created = time.Now().Add(time.Hour)
	p.Updated = time.Now().Add(-1 * time.Hour)
	p.Status = "new status"
	p.Currency = "new currency"
	change := PaymentLogChange{
		Amount:      &p.Amount,
		Description: &p.Description,
		Source:      &p.Source,
		SourceID:    &p.SourceID,
		Created:     &p.Created,
		Updated:     &p.Updated,
		Status:      &p.Status,
		Currency:    &p.Currency,
	}
	err := store.UpdatePaymentLog(p.ID, change)
	if err != nil {
		t.Errorf("Error updating payment log in memory: %s", err)
	}
	p2, ok := store.paymentLogs[p.ID]
	if !ok {
		t.Errorf("PaymentLog got lost in memory: %+v", store.paymentLogs)
	}
	success, field, expectation, result := comparePaymentLogs(p, *p2)
	if !success {
		t.Errorf("Mismatch. Expected payment log %s to be %+v, got %+v.", field, expectation, result)
	}
}

func TestUpdatingNonExistentPaymentLogInMemory(t *testing.T) {
	store := NewMemoryStore()
	newAmount := uint(100)
	err := store.UpdatePaymentLog("non-existent-payment-log", PaymentLogChange{
		Amount: &newAmount,
	})
	if err != LogNotFound {
		t.Errorf("Expected a log not found error, got %s.", err)
	}
}

func TestDeletingPaymentLogInMemory(t *testing.T) {
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
	store.paymentLogs[p.ID] = &p
	err := store.DeletePaymentLog(p.ID)
	if err != nil {
		t.Errorf("Error deleting payment log in memory: %s", err)
	}
	if _, ok := store.paymentLogs[p.ID]; ok {
		t.Errorf("Payment log was not deleted from memory as expected: %+v", store.paymentLogs)
	}
}

func TestDeletingNonExistentPaymentLogInMemory(t *testing.T) {
	store := NewMemoryStore()
	err := store.DeletePaymentLog("I don't exist")
	if err != LogNotFound {
		t.Errorf("Expected a log not found error, got %s.", err)
	}
}

func TestGettingPaymentInMemory(t *testing.T) {
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
	store.paymentLogs[p.ID] = &p
	p2, err := store.GetPaymentLog(p.ID)
	if err != nil {
		t.Errorf("Error retrieving payment log in memory: %s", err)
	}
	success, field, expectation, result := comparePaymentLogs(p, p2)
	if !success {
		t.Errorf("Mismatch: Expected payment log %s to be %+v, got %+v.", field, expectation, result)
	}
}

func TestGettingNonExistentPaymentLogInMemory(t *testing.T) {
	store := NewMemoryStore()
	_, err := store.GetPaymentLog("totally not a payment log")
	if err != LogNotFound {
		t.Errorf("Expected a log not found error, got %s.", err)
	}
}

func TestListingPaymentLogsByCampaignInMemory(t *testing.T) {
	store := NewMemoryStore()
	logs := []PaymentLog{
		PaymentLog{
			ID:          "test-payment-log 1",
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
		}, PaymentLog{
			ID:          "test-payment-log 2",
			Amount:      1,
			Source:      SourceBalanced,
			SourceID:    "balanced-id",
			Created:     time.Now().Add(time.Hour),
			Status:      StatusPending,
			Currency:    CurrencyUSD,
			CampaignID:  "campaign-id",
			GoalID:      "goal-id",
			UserID:      "user-id",
			AccountID:   "account-id",
			AccountType: "google",
		}, PaymentLog{
			ID:          "test-payment-log 3",
			Amount:      1,
			Source:      SourceBalanced,
			SourceID:    "balanced-id",
			Created:     time.Now().Add(time.Hour * 2),
			Status:      StatusPending,
			Currency:    CurrencyUSD,
			CampaignID:  "other-other-campaign-id",
			GoalID:      "goal-id",
			UserID:      "user-id",
			AccountID:   "account-id",
			AccountType: "google",
		}, PaymentLog{
			ID:          "test-payment-log 4",
			Amount:      1,
			Source:      SourceBalanced,
			SourceID:    "balanced-id",
			Created:     time.Now().Add(time.Hour * 3),
			Status:      StatusPending,
			Currency:    CurrencyUSD,
			CampaignID:  "campaign-id",
			GoalID:      "goal-id",
			UserID:      "user-id",
			AccountID:   "account-id",
			AccountType: "google",
		}, PaymentLog{
			ID:          "test-payment-log 5",
			Amount:      1,
			Source:      SourceBalanced,
			SourceID:    "balanced-id",
			Created:     time.Now().Add(time.Hour * 4),
			Status:      StatusPending,
			Currency:    CurrencyUSD,
			CampaignID:  "other-campaign-id",
			GoalID:      "goal-id",
			UserID:      "user-id",
			AccountID:   "account-id",
			AccountType: "google",
		}, PaymentLog{
			ID:          "test-payment-log 6",
			Amount:      1,
			Source:      SourceBalanced,
			SourceID:    "balanced-id",
			Created:     time.Now().Add(time.Hour * 5),
			Status:      StatusPending,
			Currency:    CurrencyUSD,
			CampaignID:  "campaign-id",
			GoalID:      "goal-id",
			UserID:      "user-id",
			AccountID:   "account-id",
			AccountType: "google",
		},
	}
	filteredLogs := make([]PaymentLog, 0)
	for pos, _ := range logs {
		if logs[pos].CampaignID == "campaign-id" {
			filteredLogs = append(filteredLogs, logs[pos])
		}
		log := logs[pos]
		store.paymentLogs[log.ID] = &log
	}
	filteredLogs = SortLogsByCreated(filteredLogs)
	results, err := store.ListPaymentLogsByCampaign("campaign-id", len(logs), 0)
	if err != nil {
		t.Errorf("Error listing payment logs by campaign: %s", err)
	}
	if len(results) != len(filteredLogs) {
		t.Logf("Log results: %+v", results)
		t.Logf("Log expectation: %+v", filteredLogs)
		t.Errorf("Expected %d payment logs, got %d.", len(filteredLogs), len(results))
	}
	for pos, _ := range results {
		success, field, expectation, result := comparePaymentLogs(filteredLogs[pos], results[pos])
		if !success {
			t.Errorf("Expected result %d %s to be %+v, got %+v.", pos, field, expectation, result)
		}
	}
}
