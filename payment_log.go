package paymentlog

import (
	"errors"
	"sort"
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
	MissingProjectID   = errors.New("Missing payment log campaign ID.")
	MissingUserID      = errors.New("Missing payment log user ID.")
	MissingAccountType = errors.New("Missing payment log account type.")
	MissingAccountID   = errors.New("Missing payment log account ID.")

	AlreadyExists = errors.New("Payment log already exists.")
	LogNotFound   = errors.New("Payment log not found.")
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
	ProjectID   string
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
	case p.ProjectID == "":
		return MissingProjectID
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
	Amount      *uint
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
	ListPaymentLogsByProject(campaignID string, num, offset int) ([]PaymentLog, error)
	ListPaymentLogsByUser(userID string, num, offset int) ([]PaymentLog, error)
	ListPaymentLogs(num, offset int) ([]PaymentLog, error)

	StoreFailureLog(failure FailureLog) error
	ListFailureLogs(num, offset int) ([]FailureLog, error)
	ListFailureLogsSince(timestamp time.Time) ([]FailureLog, error)
}

type createdSortedLogs []PaymentLog

func (c createdSortedLogs) Len() int {
	return len(c)
}

func (c createdSortedLogs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c createdSortedLogs) Less(i, j int) bool {
	return c[i].Created.After(c[j].Created)
}

func SortLogsByCreated(logs []PaymentLog) []PaymentLog {
	slogs := createdSortedLogs(logs)
	sort.Sort(slogs)
	return []PaymentLog(slogs)
}
