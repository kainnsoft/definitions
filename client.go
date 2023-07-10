package definitions

import (
	excp "gitlab.fbs-d.com/dev/go/legacy/exceptions"
	"gitlab.fbs-d.com/dev/go/legacy/helpers/api"
	"gitlab.fbs-d.com/dev/go/legacy/rmq/v2"
)

// IInternalTransfersClient интерфейс клиента для работы с сервисом Internal Transfers
type IInternalTransfersClient interface {
	Init(r rmq.IRmq)
	StartTwoPhasedInternalTransferV1(token string, client api.Client, transferKey string, fromAccountId, toAccountId, amount int64) (r *StartTwoPhasedInternalTransferResponseV1, ex excp.IException)
	CompleteTwoPhasedInternalTransferV1(token string, client api.Client, transferId int64) (r *CompleteTwoPhasedInternalTransferResponseV1, ex excp.IException)
}

// InternalTransfersClient клиент для работы с сервисом Internal Transfers
type InternalTransfersClient struct {
	Rmq rmq.IRmq
}

// Init инициализация клиента
func (c *InternalTransfersClient) Init(r rmq.IRmq) {
	c.Rmq = r
}
