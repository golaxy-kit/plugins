package nats

import (
	"github.com/nats-io/nats.go"
	"golang.org/x/net/context"
	"kit.golaxy.org/plugins/broker"
	"kit.golaxy.org/plugins/logger"
	"strings"
)

func newNatsSubscriber(ctx context.Context, nb *_NatsBroker, pattern string, opts broker.SubscriberOptions) (broker.Subscriber, error) {
	if nb.options.TopicPrefix != "" {
		pattern = nb.options.TopicPrefix + pattern
	}

	queueName := opts.QueueName
	if nb.options.QueuePrefix != "" {
		queueName = nb.options.QueuePrefix + queueName
	}

	var sub *nats.Subscription
	var err error
	var eventChan chan broker.Event
	eventHandler := opts.EventHandler

	if eventHandler == nil {
		eventChan = make(chan broker.Event, opts.EventChanSize)
	}

	ns := &_NatsSubscriber{}

	msgHandler := func(msg *nats.Msg) {
		e := _NatsEvent{
			msg: msg,
			ns:  ns,
		}

		if eventHandler != nil {
			err := eventHandler(e)
			if err != nil {
				logger.Tracef(ns.nb.ctx, "handler msg failed, %s", err)
			}
		} else {
			select {
			case eventChan <- e:
			default:
				break
			}
		}
	}

	if opts.QueueName != "" {
		sub, err = nb.client.QueueSubscribe(pattern, opts.QueueName, msgHandler)
	} else {
		sub, err = nb.client.Subscribe(pattern, msgHandler)
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	ns.cancel = cancel
	ns.nb = nb
	ns.sub = sub
	ns.options = opts
	ns.eventChan = eventChan

	go func() {
		<-ctx.Done()
		ns.Unsubscribe()
	}()

	logger.Debugf(nb.ctx, "subscribe topic %q with queue %q", pattern, queueName)

	return ns, nil
}

type _NatsSubscriber struct {
	cancel    context.CancelFunc
	nb        *_NatsBroker
	sub       *nats.Subscription
	options   broker.SubscriberOptions
	eventChan chan broker.Event
}

// Pattern returns the subscription pattern used to create the subscriber.
func (s *_NatsSubscriber) Pattern() string {
	return strings.TrimPrefix(s.sub.Subject, s.nb.options.TopicPrefix)
}

// QueueName subscribers with the same queue name will create a shared subscription where each receives a subset of messages.
func (s *_NatsSubscriber) QueueName() string {
	return strings.TrimPrefix(s.sub.Queue, s.nb.options.QueuePrefix)
}

// Unsubscribe unsubscribes the subscriber from the topic.
func (s *_NatsSubscriber) Unsubscribe() error {
	err := s.sub.Unsubscribe()
	if err != nil {
		return err
	}

	logger.Debugf(s.nb.ctx, "unsubscribe topic %q with %q", s.sub.Subject, s.sub.Queue)

	s.cancel()
	return nil
}

// Next is a blocking call that waits for the next event to be received from the subscriber.
func (s *_NatsSubscriber) Next() (broker.Event, error) {
	for event := range s.eventChan {
		return event, nil
	}
	return nil, broker.ErrUnsubscribed
}