package service

import (
	"context"

	"github.com/bitcapybara/geckod"
)

type Subscription interface {
	GetName() string
	GetType() geckod.SubScriptionType

	AddConsumer(ctx context.Context, _ *Consumer) error
	DelConsumer(ctx context.Context, consumerId uint64) error

	Flow(ctx context.Context, consumerId uint64, permits uint64) error

	Ack(ctx context.Context, ackType geckod.AckType, msgIds []uint64) error
	Unsubscribe(ctx context.Context, consumer *Consumer) error

	// 有消费者，才会生成 dispatcher，否则返回 NotFound
	GetDispatcher() (Dispatcher, error)

	Close(ctx context.Context) error
}

type Dispatcher interface {
	AddConsumer(ctx context.Context, _ *Consumer) error
	DelConsumer(ctx context.Context, consumerId uint64) error
	GetConsumers(ctx context.Context) []*Consumer

	Flow(ctx context.Context, consumerId uint64, permits uint64) error
	CanUnsubscribe(ctx context.Context, consumerId uint64) bool
	SendMessages(ctx context.Context) error

	Close(ctx context.Context) error
}

type Cursor interface {
	// id之前的所有记录都删除
	MarkDelete(ctx context.Context, msgId uint64) error
	// 删除所有指定 id 的记录
	Delete(ctx context.Context, msgIds []uint64) error
	// 根据cursor的情况，获取可发送的消息
	GetMoreMessages(ctx context.Context, num uint64) ([]*geckod.RawMessage, error)
}
