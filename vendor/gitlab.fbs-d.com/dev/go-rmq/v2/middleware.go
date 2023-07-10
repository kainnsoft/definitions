package rmq

import (
	"context"

	"github.com/streadway/amqp"
)

type (
	ListenerType int8
	ContextType  int8

	Next       func() error
	Middleware func(*MiddlewareContext) error

	MiddlewareContext struct {
		ListenerType ListenerType
		ContextType  ContextType

		index         int
		ListenerRPC   func(*MiddlewareContext) (interface{}, error)
		ListenerEvent func(*MiddlewareContext) error
		Middleware    []Middleware

		MethodName  string
		ServiceName string

		AmqpDelivery   *amqp.Delivery
		AmqpPublishing *amqp.Publishing
		Request        interface{}
		GetRequest     func() interface{}
		Response       interface{}
		Context        context.Context

		AmqpDeliveries []*amqp.Delivery
		Requests       []interface{}
	}
)

const (
	RPC    ListenerType = iota
	Event               // single message
	Events              // batch messages

	Client ContextType = iota
	Server
)

func (ctx *MiddlewareContext) Next() (err error) {
	if ctx.index < len(ctx.Middleware) {
		ctx.index++
		return ctx.Middleware[ctx.index-1](ctx)
	}

	if ctx.ListenerType == RPC {
		if ctx.Response, err = ctx.ListenerRPC(ctx); err != nil {
			return err
		}
		return nil
	}

	return ctx.ListenerEvent(ctx)
}
