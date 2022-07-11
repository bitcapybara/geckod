package service

import (
	"github.com/bitcapybara/geckod"
)

type Subscription interface {
	GetName() string
	GetType() geckod.SubScriptionType

	AddConsumer(*Consumer) error
	DelConsumer(consumerId uint64) error

	Flow(consumerId uint64, permits uint64) error

	Ack(ackType geckod.AckType, msgIds []uint64) error
	Unsubscribe(consumer *Consumer) error

	// 有消费者，才会生成 dispatcher，否则返回 NotFound
	GetDispatcher() (Dispatcher, error)

	Close() error
}

type Dispatcher interface {
	AddConsumer(*Consumer) error
	DelConsumer(consumerId uint64) error
	GetConsumers() []*Consumer

	Flow(consumerId uint64, permits uint64) error
	CanUnsubscribe(consumerId uint64) bool
	SendMessages() error

	Close() error
}

type Cursor interface {
	// id之前的所有记录都删除
	MarkDelete(msgId uint64) error
	// 删除所有指定 id 的记录
	Delete(msgIds []uint64) error
}
