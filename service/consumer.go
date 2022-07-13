package service

import (
	"context"
	"sync"

	"github.com/bitcapybara/geckod"
	cmdpb "github.com/bitcapybara/geckod-proto/gen/go/proto/command"
	"github.com/bitcapybara/geckod/errs"

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

	mu sync.Mutex
	// 仅共享订阅模式消费者使用
	// 某一消息交给当前消费者发送后，后续不会再给其他消费者重试
	// 重发机制：
	// 1. 将消息置为 nack
	// 2. 等待 ack 超时
	// 3. 死信队列
	pendingAcks map[uint64]struct{}
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
	// 累积ack不可用于共享订阅
	if ackType == geckod.AckType(cmdpb.Ack_Cumulative) {
		if c.sub.GetType() == geckod.SubScriptionType(cmdpb.Subscribe_Shared) {
			return errs.ErrUnmatchAckType
		}
	}

	if err := c.sub.Ack(ctx, ackType, msgIds); err != nil {
		return err
	}

	// 单个ack，处理 pendingack
	// pendingack 只有共享订阅才有，所以一定不会是累积ack
	if ackType == geckod.AckType(cmdpb.Ack_Individual) {
		c.removePendingAcks(msgIds)
	}

	return nil
}

func (c *Consumer) removePendingAcks(msgIds []uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, id := range msgIds {
		delete(c.pendingAcks, id)
	}
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
