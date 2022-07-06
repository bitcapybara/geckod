package service

import (
	"github.com/bitcapybara/geckod"
)

type Subscription interface {
	GetName() string
	GetType() geckod.SubScriptionType

	AddConsumer(Consumer) error
	DelConsumer(consumerId uint64) error

	Flow(consumerId uint64, permits uint64) error

	Ack(ackType geckod.AckType, msgIds []uint64) error
	Unsubscribe(consumerId uint64) error

	Close() error
}

type Dispatcher interface {
	AddConsumer(Consumer) error
	DelConsumer(consumerId uint64) error

	Flow(consumerId uint64, permits uint64) error

	Close() error
}
