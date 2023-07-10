// Code generated by gitlab.fbs-d.com/dev/amqp-rpc-generator/protoc-gen-amqprpc3 v1.3.1 DO NOT EDIT.
// source: proto/common_partner-transfers.rmq.client.proto

package partner_transfers_defs

import (
	"context"
	"fmt"
	"strings"

	"github.com/streadway/amqp"
	"gitlab.fbs-d.com/dev/go-rmq/v2"
)

var (
	_ = fmt.Errorf
	_ = strings.Join
	_ = amqp.Dial
)

// Server API
type RMQClient interface {
	StartTwoPhasedPartnerTransferV1(context.Context, *StartTwoPhasedPartnerTransferV1Request) (*StartTwoPhasedPartnerTransferV1Response, error)
	CompleteTwoPhasedPartnerTransferV1(context.Context, *CompleteTwoPhasedPartnerTransferV1Request) (*CompleteTwoPhasedPartnerTransferV1Response, error)
}
type serverRMQClient struct {
	rmq          rmq.Rmq
	instanceName string
	handlers     RMQClient
	runners      map[string]func() error
	opt          []rmq.Middleware
}

func NewServerRMQClient(
	instanceName string,
	rmq rmq.Rmq,
	handlers RMQClient,
	opt []rmq.Middleware,
) *serverRMQClient {
	var s = &serverRMQClient{
		instanceName: instanceName,
		rmq:          rmq,
		handlers:     handlers,
		opt:          opt,
	}
	s.register()
	return s
}

// RMQ receiver
func (r *serverRMQClient) RMQ() rmq.Rmq {
	return r.rmq
}

// Run all
func (r *serverRMQClient) Run() error {
	go r.RunByName("StartTwoPhasedPartnerTransferV1")
	go r.RunByName("CompleteTwoPhasedPartnerTransferV1")
	return nil
}

// RunByName
func (r *serverRMQClient) RunByName(name string) error {
	if runner, ok := r.runners[name]; ok {
		go func() {
			if err := runner(); err != nil {
				panic(fmt.Sprintln("rmq init RMQClient runner", name, err))
			}
		}()
		return nil
	}
	return fmt.Errorf("not found")
}

// register functions
func (r *serverRMQClient) register() {
	r.runners = map[string]func() error{
		"StartTwoPhasedPartnerTransferV1": func() error {
			return r.rmq.ConsumeRpc("r.partner-transfers.StartTwoPhasedPartnerTransfer.v1", _Handler_StartTwoPhasedPartnerTransferV1(r.opt, r.handlers.StartTwoPhasedPartnerTransferV1))
		},
		"CompleteTwoPhasedPartnerTransferV1": func() error {
			return r.rmq.ConsumeRpc("r.partner-transfers.CompleteTwoPhasedPartnerTransfer.v1", _Handler_CompleteTwoPhasedPartnerTransferV1(r.opt, r.handlers.CompleteTwoPhasedPartnerTransferV1))
		},
	}
}

// Server API handlers
func _Handler_StartTwoPhasedPartnerTransferV1(
	opt []rmq.Middleware,
	h func(context.Context, *StartTwoPhasedPartnerTransferV1Request) (*StartTwoPhasedPartnerTransferV1Response, error),
) func(*amqp.Delivery) *amqp.Publishing {
	return func(del *amqp.Delivery) *amqp.Publishing {
		var mdwCtx = &rmq.MiddlewareContext{
			ListenerRPC: func(ctx *rmq.MiddlewareContext) (interface{}, error) {
				return h(ctx.Context, ctx.Request.(*StartTwoPhasedPartnerTransferV1Request))
			},
			ListenerType: rmq.RPC,
			ContextType:  rmq.Server,
			Middleware:   opt,
			MethodName:   "StartTwoPhasedPartnerTransferV1",
			ServiceName:  "RMQClient",
			AmqpDelivery: del,
			Request:      &StartTwoPhasedPartnerTransferV1Request{},
			Response:     &StartTwoPhasedPartnerTransferV1Response{},
			Context:      context.Background(),
		}

		_ = mdwCtx.Next()
		return mdwCtx.AmqpPublishing
	}
}
func _Handler_CompleteTwoPhasedPartnerTransferV1(
	opt []rmq.Middleware,
	h func(context.Context, *CompleteTwoPhasedPartnerTransferV1Request) (*CompleteTwoPhasedPartnerTransferV1Response, error),
) func(*amqp.Delivery) *amqp.Publishing {
	return func(del *amqp.Delivery) *amqp.Publishing {
		var mdwCtx = &rmq.MiddlewareContext{
			ListenerRPC: func(ctx *rmq.MiddlewareContext) (interface{}, error) {
				return h(ctx.Context, ctx.Request.(*CompleteTwoPhasedPartnerTransferV1Request))
			},
			ListenerType: rmq.RPC,
			ContextType:  rmq.Server,
			Middleware:   opt,
			MethodName:   "CompleteTwoPhasedPartnerTransferV1",
			ServiceName:  "RMQClient",
			AmqpDelivery: del,
			Request:      &CompleteTwoPhasedPartnerTransferV1Request{},
			Response:     &CompleteTwoPhasedPartnerTransferV1Response{},
			Context:      context.Background(),
		}

		_ = mdwCtx.Next()
		return mdwCtx.AmqpPublishing
	}
}

// Client API handlers
type clientRMQClient struct {
	rmq rmq.Rmq
	opt []rmq.Middleware
}

func NewClientRMQClient(rmq rmq.Rmq, opt []rmq.Middleware) *clientRMQClient {
	return &clientRMQClient{
		rmq: rmq,
		opt: opt,
	}
}

// Client API handlers
func (r *clientRMQClient) StartTwoPhasedPartnerTransferV1(
	ctx context.Context,
	arg *StartTwoPhasedPartnerTransferV1Request,
) (out *StartTwoPhasedPartnerTransferV1Response, err error) {
	var mdwCtx = &rmq.MiddlewareContext{
		ListenerRPC: func(ctx *rmq.MiddlewareContext) (_ interface{}, err error) {
			ctx.AmqpDelivery, err = r.rmq.Rpc(ctx.Context, "r.partner-transfers.StartTwoPhasedPartnerTransfer.v1", ctx.AmqpPublishing)
			return ctx.Response, err
		},
		Request:      arg,
		Middleware:   r.opt,
		ListenerType: rmq.RPC,
		ContextType:  rmq.Client,
		MethodName:   "StartTwoPhasedPartnerTransferV1",
		ServiceName:  "RMQClient",
		Response:     &StartTwoPhasedPartnerTransferV1Response{},
		Context:      ctx,
	}

	if err = mdwCtx.Next(); err != nil {
		return nil, err
	}
	return mdwCtx.Response.(*StartTwoPhasedPartnerTransferV1Response), nil
}
func (r *clientRMQClient) CompleteTwoPhasedPartnerTransferV1(
	ctx context.Context,
	arg *CompleteTwoPhasedPartnerTransferV1Request,
) (out *CompleteTwoPhasedPartnerTransferV1Response, err error) {
	var mdwCtx = &rmq.MiddlewareContext{
		ListenerRPC: func(ctx *rmq.MiddlewareContext) (_ interface{}, err error) {
			ctx.AmqpDelivery, err = r.rmq.Rpc(ctx.Context, "r.partner-transfers.CompleteTwoPhasedPartnerTransfer.v1", ctx.AmqpPublishing)
			return ctx.Response, err
		},
		Request:      arg,
		Middleware:   r.opt,
		ListenerType: rmq.RPC,
		ContextType:  rmq.Client,
		MethodName:   "CompleteTwoPhasedPartnerTransferV1",
		ServiceName:  "RMQClient",
		Response:     &CompleteTwoPhasedPartnerTransferV1Response{},
		Context:      ctx,
	}

	if err = mdwCtx.Next(); err != nil {
		return nil, err
	}
	return mdwCtx.Response.(*CompleteTwoPhasedPartnerTransferV1Response), nil
}
