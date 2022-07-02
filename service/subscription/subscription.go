package subscription

import "github.com/bitcapybara/geckod/service/consumer"

type Subscription interface {
	GetName() string

	GetDispatcher() Dispatcher

	Ack(msgId uint64) error
	Unsubscribe(consumerId uint64) error

	Close() error
}

type Dispatcher interface {
	AddConsumer(consumer.Consumer) error
	DelConsumer(consumerId uint64) error

	Flow(consumerId uint64, permits uint64) error

	Close() error
}
