package definitions

import (
	"encoding/json"
	"time"

	"gitlab.fbs-d.com/dev/go/legacy/helpers/api"
)

type CreatedInternalTransferEvent struct {
	Token  string                            `json:"token"`
	Client api.Client                        `json:"client"`
	Body   *CreatedInternalTransferEventBody `json:"body"`
}

type CreatedInternalTransferEventBody struct {
	UserId              int64  `json:"userId"`
	AccountId           int64  `json:"accountId"`
	AccountIdExternal   int64  `json:"accountIdExternal"`
	TransactionComment  string `json:"transactionComment"`
	TransactionAmount   int64  `json:"transactionAmount"`
	TransactionCurrency string `json:"transactionCurrency"`
	TransactionId       int64  `json:"transactionId"`
}

func (res *CreatedInternalTransferEvent) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}

type CompletedInternalTransferEvent struct {
	Token  string                              `json:"token"`
	Client api.Client                          `json:"client"`
	Body   *CompletedInternalTransferEventBody `json:"body"`
}

type CompletedInternalTransferEventBody struct {
	UserId              int64     `json:"userId"`
	AccountId           int64     `json:"accountId"`
	AccountIdExternal   int64     `json:"accountIdExternal"`
	TransactionComment  string    `json:"transactionComment"`
	TransactionAmount   int64     `json:"transactionAmount"`
	TransactionCurrency string    `json:"transactionCurrency"`
	TransactionId       int64     `json:"transactionId"`
	CreatedAt           time.Time `json:"created_at"`
	AcceptedAt          time.Time `json:"accepted_at"`
}

func (res *CompletedInternalTransferEvent) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}
