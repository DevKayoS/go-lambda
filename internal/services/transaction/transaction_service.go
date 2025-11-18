package transaction

import "github.com/DevKayoS/go-lambda/internal/models"

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) Create(body models.TransactionRequest) (any, error) {
	return body, nil
}
