package impl

import (
	"context"
	"errors"
	"sync"

	"github.com/bitcapybara/geckod"
	cmdpb "github.com/bitcapybara/geckod-proto/gen/go/proto/command"
	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
)

var _ service.Topic = (*topic)(nil)

type topic struct {
	name string

	consumers service.Consumers
	storage   service.Storage

	mu                   sync.Mutex
	producers            map[uint64]*service.Producer
	subscriptions        map[string]service.Subscription
	exclusiveProducer    uint64
	hasExclusiveProducer bool
}

func newTopic(name string) *topic {
	return nil
}

func (t *topic) GetName() string {
	return t.name
}

// 处理客户端生产者发送的数据
func (t *topic) Publish(ctx context.Context, msg *geckod.RawMessage) error {
	// 存储
	if _, err := t.storage.Add(ctx, msg); err != nil {
		return err
	}
	//
	for _, subscripton := range t.subscriptions {
		dispatcher, err := subscripton.GetDispatcher()
		if err != nil {
			if err == errs.ErrNotFound {
				continue
			}
			return err
		}
		if err := dispatcher.SendMessages(ctx); err != nil {
			return err
		}
	}
	return nil
}

// 处理消费者订阅
// 生成 consumer，添加到 subscription，返回 consumer
func (t *topic) Subscribe(ctx context.Context, option *service.SubscriptionOption) (*service.Consumer, error) {

	if len(option.SubName) == 0 {
		return nil, errors.New("subscription name is empty")
	}

	subscription := t.getOrCreateSubscription(option.SubName)

	consumer, err := t.consumers.GetOrCreate(&service.AddConsumerParams{
		ClientId:     option.ClientId,
		ConsumerName: option.ConsumerName,
		TopicName:    option.TopicName,
		Subscription: subscription,
	})
	if err != nil {
		return nil, err
	}

	if err := subscription.AddConsumer(ctx, consumer); err != nil {
		return nil, err
	}

	return consumer, nil
}

func (t *topic) Unsubscribe(ctx context.Context, subName string) error {
	return nil
}

func (t *topic) RemoveSubscription(ctx context.Context, subName string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.subscriptions, subName)
	return nil
}

func (t *topic) getOrCreateSubscription(name string) service.Subscription {
	t.mu.Lock()
	defer t.mu.Unlock()

	if sub, ok := t.subscriptions[name]; ok {
		return sub
	}
	sub := newSubscription()
	t.subscriptions[name] = sub
	return sub
}

// 生产者管理
func (t *topic) AddProducer(ctx context.Context, p *service.Producer) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 如果已经有了 exclusive 的生产者，则其他生产者不允许添加
	if p.AccessMode == geckod.ProducerAccessMode(cmdpb.Producer_Exclusive) {
		if t.hasExclusiveProducer {
			return errs.ErrProducerExclusiveConflict
		}
	}

	if _, ok := t.producers[p.Id]; ok {
		return errs.ErrProducerAlreadyExists
	}

	t.producers[p.Id] = p
	return nil
}

func (t *topic) GetProducer(ctx context.Context, producer_id uint64) (*service.Producer, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if p, ok := t.producers[producer_id]; ok {
		return p, nil
	}
	return nil, errs.ErrNotFound
}

func (t *topic) DelProducer(ctx context.Context, producer_id uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.producers, producer_id)
	return nil
}

// 释放资源
func (t *topic) Close(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}
