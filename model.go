package rinha

import (
	"time"

	"github.com/jellydator/validation"
)

type TransactionType string

const (
	CreditType TransactionType = "c"
	DebitType  TransactionType = "d"
)

type Transaction struct {
	Amount      uint            `json:"valor"`
	Type        TransactionType `json:"tipo"`
	Description string          `json:"descricao"`
	CreatedAt   time.Time       `json:"realizada_em"`
}

func (t Transaction) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Amount, validation.Required),
		validation.Field(&t.Type, validation.Required, validation.In(CreditType, DebitType)),
		validation.Field(&t.Description, validation.Required, validation.Length(1, 10)),
	)
}

type Client struct {
	ID             uint          `db:"id" json:"id"`
	AccountLimit   int           `db:"limite" json:"limite"`
	AccountBalance int           `db:"saldo" json:"saldo"`
	Transactions   []Transaction `db:"ultimas_transacoes" json:"ultimas_transacoes"`
}

type Balance struct {
	AccountLimit   int `db:"limite" json:"limite"`
	AccountBalance int `db:"saldo" json:"saldo"`
}
