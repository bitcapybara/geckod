package service

import (
	"context"
	"sync"

	"github.com/bitcapybara/geckod"
	"go.uber.org/atomic"
)

type Consumers interface {
	GetOrCreate(*AddConsumerParams) (*Consumer, error)
	Get(id uint64) (*Consumer, error)
	Del(id uint64)
	Add(*Consumer) error
}

type AddConsumerParams struct {
	ClientId     uint64
	ConsumerName string
	TopicName    string
	Subscription Subscription
}

type Consumer struct {
	Id        uint64
	Name      string
	ClientId  uint64
	TopicName string

	sub    Subscription
	sender geckod.ConsumerMessageSender

	permits atomic.Uint64

	mu          sync.Mutex
	pendingAcks []uint64
}

func NewConsumer(id uint64, params *AddConsumerParams) *Consumer {
	return &Consumer{
		Id:        id,
		Name:      params.ConsumerName,
		ClientId:  params.ClientId,
		TopicName: params.TopicName,
		sub:       params.Subscription,
	}
}

func (c *Consumer) Unsubscribe(ctx context.Context) error {
	return c.sub.Unsubscribe(ctx, c)
}

func (c *Consumer) Flow(ctx context.Context, permits uint64) error {
	defer c.permits.Add(permits)
	return c.sub.Flow(ctx, c.Id, permits)
}

func (c *Consumer) Ack(ctx context.Context, ackType geckod.AckType, msgIds []uint64) error {
	if err := c.sub.GetType().MatchAckType(ackType); err != nil {
		return err
	}
	return c.sub.Ack(ctx, ackType, msgIds)
}

// 把消息发送给消费者
func (c *Consumer) SendMessages(ctx context.Context, msgs []*geckod.RawMessage) error {
	defer c.permits.Add(-uint64(len(msgs)))
	return nil
}

func (c *Consumer) Permits() uint64 {
	return c.permits.Load()
}

func (c *Consumer) Close(ctx context.Context) error {
	return c.sub.DelConsumer(ctx, c.Id)
}
