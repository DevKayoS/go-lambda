package models

type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "debit"
	TransactionTypeCredit TransactionType = "credit"
)

type TransactionRequest struct {
	Amount      int64           `json:"amount"`
	Description string          `json:"description"`
	Type        TransactionType `json:"type"`
}
