package definitions

import (
	"encoding/json"

	excp "gitlab.fbs-d.com/dev/go/legacy/exceptions"
	toolkit "gitlab.fbs-d.com/dev/go/legacy/helpers"
	"gitlab.fbs-d.com/dev/go/legacy/helpers/api"
)

const (
	RPCStartTwoPhasedInternalTransferV1 = "r.internal-transfers.StartTwoPhasedInternalTransfer.v1"
)

// Запрос на внутренний перевод средств между счетами
type StartTwoPhasedInternalTransferRequestV1 struct {
	Token  string                                      `json:"token"`
	Client api.Client                                  `json:"client"`
	Body   StartTwoPhasedInternalTransferRequestV1Body `json:"body"`
}

type StartTwoPhasedInternalTransferRequestV1Body struct {
	TransferKey   string `json:"transferKey"`
	FromAccountId int64  `json:"fromAccountId" validate:"positive"`
	ToAccountId   int64  `json:"toAccountId" validate:"positive"`
	Amount        int64  `json:"amount" validate:"positive"`
}

func (res *StartTwoPhasedInternalTransferRequestV1) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}

type StartTwoPhasedInternalTransferResponseV1 struct {
	OperationId int64 `json:"operationId"`
}

func (res *StartTwoPhasedInternalTransferResponseV1) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}

// StartTwoPhasedInternalTransferV1 Создает внутренний перевод средств между счетами
func (c *InternalTransfersClient) StartTwoPhasedInternalTransferV1(token string, client api.Client, transferKey string, fromAccountId, toAccountId, amount int64) (r *StartTwoPhasedInternalTransferResponseV1, ex excp.IException) {

	rq := StartTwoPhasedInternalTransferRequestV1{
		Token:  token,
		Client: client,
		Body: StartTwoPhasedInternalTransferRequestV1Body{
			TransferKey:   transferKey,
			Amount:        amount,
			FromAccountId: fromAccountId,
			ToAccountId:   toAccountId,
		},
	}

	uuid, ex := toolkit.Uuid()
	if ex != nil {
		return
	}

	res, ex := c.Rmq.Rpc(RPCStartTwoPhasedInternalTransferV1, uuid, rq.GetJson())
	if ex != nil {
		return
	}

	r = new(StartTwoPhasedInternalTransferResponseV1)
	err := json.Unmarshal(res, r)
	if err != nil {
		ex = excp.NewJsonUnmarshalException(err)
		return
	}
	return

}
