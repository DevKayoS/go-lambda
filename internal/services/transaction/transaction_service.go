package transaction

import "github.com/DevKayoS/go-lambda/internal/models"

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

// TODO: implementar o servico de transacao com create, read, update
func (t *TransactionService) Create(body models.TransactionRequest) (any, error) {
	return body, nil
}
