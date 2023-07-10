package definitions

import (
	"encoding/json"

	excp "gitlab.fbs-d.com/dev/go/legacy/exceptions"
	toolkit "gitlab.fbs-d.com/dev/go/legacy/helpers"
	"gitlab.fbs-d.com/dev/go/legacy/helpers/api"
)

const (
	RPCCompleteTwoPhasedInternalTransferV1 = "r.partner-transfers.CompleteTwoPhasedInternalTransfer.v1"
)

// Запрос на завершение внутреннего перевода средств между счетами
type CompleteTwoPhasedInternalTransferRequestV1 struct {
	Token  string                                         `json:"token"`
	Client api.Client                                     `json:"client"`
	Body   CompleteTwoPhasedInternalTransferRequestV1Body `json:"body"`
}

type CompleteTwoPhasedInternalTransferRequestV1Body struct {
	OperationId int64 `json:"operationId" validate:"positive"`
}

func (res *CompleteTwoPhasedInternalTransferRequestV1) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}

type CompleteTwoPhasedInternalTransferResponseV1 struct {
}

func (res *CompleteTwoPhasedInternalTransferResponseV1) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}

// CompleteTwoPhasedInternalTransferV1 Завершает внутренний перевод средств между счетами
func (c *InternalTransfersClient) CompleteTwoPhasedInternalTransferV1(token string, client api.Client, transferId int64) (r *CompleteTwoPhasedInternalTransferResponseV1, ex excp.IException) {

	rq := CompleteTwoPhasedInternalTransferRequestV1{
		Token:  token,
		Client: client,
		Body: CompleteTwoPhasedInternalTransferRequestV1Body{
			OperationId: transferId,
		},
	}

	uuid, ex := toolkit.Uuid()
	if ex != nil {
		return
	}

	res, ex := c.Rmq.Rpc(RPCCompleteTwoPhasedInternalTransferV1, uuid, rq.GetJson())
	if ex != nil {
		return
	}

	r = new(CompleteTwoPhasedInternalTransferResponseV1)
	err := json.Unmarshal(res, r)
	if err != nil {
		ex = excp.NewJsonUnmarshalException(err)
		return
	}
	return

}
