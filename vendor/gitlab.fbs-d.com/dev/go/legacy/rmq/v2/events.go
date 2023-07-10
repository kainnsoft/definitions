package rmq

import (
	"encoding/json"
	"errors"
	"gitlab.fbs-d.com/dev/go/legacy/exceptions"
	"log"
	"gitlab.fbs-d.com/dev/go/legacy/helpers/api"
	"time"
)

type EventMessage struct {
	Name  string
	Event EventMessagePayload
}

type EventMessagePayload struct {
	Token  string      `json:"token"`
	Client api.Client  `json:"client"`
	Body   interface{} `json:"body"`
}

func (em *EventMessage) GetJson() (data []byte) {
	data, _ = json.Marshal(em.Event)
	return
}

//go:generate mockgen -source events.go -destination events_mocks.go  -self_package . -package rmq -imports "excp=gitlab.fbs-d.com/dev/go/legacy/exceptions"
type IEventSender interface {
	Run(r IRmq) error
	Trigger(token string, client api.Client, eventName string, event interface{}) excp.IException
}

type EventSender struct {
	Rmq            IRmq
	Exchange       string
	RepeatInterval time.Duration
	ch             chan EventMessage
}

func (es *EventSender) Run(r IRmq) error {
	if es.Exchange == "" {
		return errors.New("Empty exchange name")
	}

	es.Rmq = r
	es.ch = make(chan EventMessage)
	if es.RepeatInterval == 0 {
		es.RepeatInterval = 5 * time.Second
	}

	go es.Listen()

	return nil
}

func (es *EventSender) Listen() {
	var e EventMessage
	var ex excp.IException

	for {
		e = <-es.ch

		log.Printf("%s Trying to send %#v", api.SysToken, string(es.Rmq.TrimFieldsForLog(e.GetJson())))
		for {
			ex = es.Rmq.ExchangePublish(es.Exchange, "direct", e.Name, e.GetJson())
			if ex == nil {
				log.Printf("%s Event %s sent", e.Event.Token, e.Name)
				break
			}

			if ex.Code() == RabbitInvalidParamException {
				log.Printf("%s Invalid exchange params. Event discarded.", e.Event.Token)
				break
			}

			log.Printf("%s Failed to send %s event: %v", e.Event.Token, e.Name, ex)
			time.Sleep(es.RepeatInterval)
		}
	}
}

func (es *EventSender) Trigger(token string, client api.Client, eventName string, event interface{}) (ex excp.IException) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s Error queueing event %s: %v", token, eventName, r)
			ex = NewRabbitTriggerEventException(errors.New("Recovered from panic"))
		}
	}()

	em := EventMessage{
		Name: eventName,
		Event: EventMessagePayload{
			Token:  token,
			Client: client,
			Body:   event,
		},
	}

	es.ch <- em

	return
}
