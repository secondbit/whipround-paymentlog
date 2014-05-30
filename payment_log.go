package paymentlog

import (
	"errors"
	"time"
)

const (
	SourceBalanced = "balanced"
	StatusPending  = "pending"
	CurrencyUSD    = "usd"
)

var (
	MissingID          = errors.New("Missing payment log ID.")
	MissingAmount      = errors.New("Missing payment log amount.")
	MissingSource      = errors.New("Missing payment log source.")
	MissingSourceID    = errors.New("Missing payment log source ID.")
	MissingCreated     = errors.New("Missing payment log created timestamp.")
	MissingStatus      = errors.New("Missing payment log status.")
	MissingCurrency    = errors.New("Missing payment log currency.")
	MissingCampaignID  = errors.New("Missing payment log campaign ID.")
	MissingGoalID      = errors.New("Missing payment log goal ID.")
	MissingUserID      = errors.New("Missing payment log user ID.")
	MissingAccountType = errors.New("Missing payment log account type.")
	MissingAccountID   = errors.New("Missing payment log account ID.")

	AlreadyExists = errors.New("Payment log already exists.")
)

type PaymentLog struct {
	ID          string
	Amount      uint
	Description string
	Source      string
	SourceID    string
	Created     time.Time
	Updated     time.Time
	Status      string
	Currency    string
	CampaignID  string
	GoalID      string
	UserID      string
	AccountID   string
	AccountType string
}

func (p PaymentLog) Validate() error {
	switch {
	case p.ID == "":
		return MissingID
	case p.Amount == 0:
		return MissingAmount
	case p.Source == "":
		return MissingSource
	case p.SourceID == "":
		return MissingSourceID
	case p.Created.IsZero():
		return MissingCreated
	case p.Status == "":
		return MissingStatus
	case p.Currency == "":
		return MissingCurrency
	case p.CampaignID == "":
		return MissingCampaignID
	case p.GoalID == "":
		return MissingGoalID
	case p.UserID == "":
		return MissingUserID
	case p.AccountType == "":
		return MissingAccountType
	case p.AccountID == "":
		return MissingAccountID
	default:
		return nil
	}
}

type PaymentLogChange struct {
	Amount      *int
	Description *string
	Source      *string
	SourceID    *string
	Created     *time.Time
	Updated     *time.Time
	Status      *string
	Currency    *string
}

type FailureLog struct {
	ID                string
	PaymentLogID      string
	FailureReason     string
	FailureReasonCode string
	Timestamp         time.Time
}

type LogStore interface {
	StorePaymentLog(log PaymentLog) error
	UpdatePaymentLog(id string, change PaymentLogChange) error
	DeletePaymentLog(id string) error
	GetPaymentLog(id string) (PaymentLog, error)
	ListPaymentLogsByCampaign(campaignID string, num, offset int) ([]PaymentLog, error)
	ListPaymentLogsByGoal(campaignID, goalID string, num, offset int) ([]PaymentLog, error)
	ListPaymentLogsByUser(userID string, num, offset int) ([]PaymentLog, error)
	ListPaymentLogs(num, offset int) ([]PaymentLog, error)

	StoreFailureLog(failure FailureLog) error
	ListFailureLogs(num, offset int) ([]FailureLog, error)
	ListFailureLogsSince(timestamp time.Time) ([]FailureLog, error)
}
