package service

import (
	"context"

	"github.com/bitcapybara/geckod"
)

type Topics interface {
	GetOrCreate(name string) Topic
	Get(name string) (Topic, error)
	Del(name string)
}

type Topic interface {
	GetName() string

	// 处理客户端生产者发送的数据
	Publish(ctx context.Context, msg *geckod.RawMessage) error
	// 处理消费者订阅
	// 生成 consumer，添加到 subscription，返回 consumer
	Subscribe(ctx context.Context, subOpt *SubscriptionOption) (*Consumer, error)
	Unsubscribe(ctx context.Context, subName string) error
	RemoveSubscription(ctx context.Context, subName string) error

	// 生产者管理
	AddProducer(ctx context.Context, _ *Producer) error
	GetProducer(ctx context.Context, producer_id uint64) (*Producer, error)
	DelProducer(ctx context.Context, producer_id uint64) error

	// 释放资源
	Close(ctx context.Context) error
}

type SubscriptionOption struct {
	SubName      string
	ClientId     uint64
	ConsumerName string
	TopicName    string
}
