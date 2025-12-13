package transaction

import (
	"testing"

	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction_CreateSuccess(t *testing.T) {
	service := &TransactionService{}

	body := models.TransactionRequest{
		Amount:      123,
		Description: "transacao teste",
		Type:        "debit",
	}

	_, err := service.Create(body)

	assert.NoError(t, err)
}
