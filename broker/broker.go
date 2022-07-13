package broker

import (
	"context"
	"time"

	"github.com/bitcapybara/geckod"
	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
)

type Broker interface {
	Connect(ctx context.Context, _ *geckod.CommandConnect) (*geckod.CommandConnected, error)
	Producer(ctx context.Context, _ *geckod.CommandProducer) (*geckod.CommandProducerSuccess, error)
	Subscribe(ctx context.Context, _ *geckod.CommandSubscribe) (*geckod.CommandSubscribeSuccess, error)
	Unsubscribe(ctx context.Context, _ *geckod.CommandUnsubscribe) error
	Flow(ctx context.Context, _ *geckod.CommandFlow) error
	Ack(ctx context.Context, _ *geckod.CommandAck) error
	Send(ctx context.Context, _ *geckod.CommandSend) error
}

type Authenticator = func(username, passwd string, method geckod.ConnectAuthMethod) (bool, error)

var _ Broker = (*broker)(nil)

type broker struct {
	topics    service.Topics
	clients   service.Clients
	producers service.Producers
	consumers service.Consumers

	authFn Authenticator
}

func (b *broker) Connect(ctx context.Context, cmd *geckod.CommandConnect) (*geckod.CommandConnected, error) {
	ok, err := b.authFn(cmd.Username, cmd.Password, geckod.ConnectAuthMethod(cmd.AuthMethod))
	if !ok {
		return nil, errs.ErrAuthFailed
	}
	if err != nil {
		return nil, err
	}

	client_id, err := b.clients.Add(cmd.ClientName)
	if err != nil {
		return nil, err
	}

	return &geckod.CommandConnected{
		ClientId: client_id,
	}, nil
}

func (b *broker) Producer(ctx context.Context, cmd *geckod.CommandProducer) (*geckod.CommandProducerSuccess, error) {
	topic := b.topics.GetOrCreate(cmd.Topic)

	// 生成 producer
	producer, err := b.producers.Create(&service.AddProducerParams{
		ProducerName: cmd.ProducerName,
		AccessMode:   geckod.ProducerAccessMode(cmd.AccessMode),
		Topic:        topic,
	})
	if err != nil {
		return nil, err
	}

	// 添加到 topic
	if err := topic.AddProducer(ctx, producer); err != nil {
		return nil, err
	}

	return &geckod.CommandProducerSuccess{
		ProducerId: producer.Id,
	}, nil
}

func (b *broker) Subscribe(ctx context.Context, cmd *geckod.CommandSubscribe) (*geckod.CommandSubscribeSuccess, error) {
	topic := b.topics.GetOrCreate(cmd.Topic)

	consumer, err := topic.Subscribe(ctx, &service.SubscriptionOption{
		SubName:      cmd.SubName,
		ClientId:     cmd.ClientId,
		ConsumerName: cmd.ConsumerName,
		TopicName:    cmd.Topic,
	})
	if err != nil {
		return nil, err
	}

	if err := b.consumers.Add(consumer); err != nil {
		return nil, err
	}

	return &geckod.CommandSubscribeSuccess{
		ConsumerId: consumer.Id,
	}, nil
}

func (b *broker) Unsubscribe(ctx context.Context, cmd *geckod.CommandUnsubscribe) error {

	consumer, err := b.consumers.Get(cmd.ConsumerId)
	if err != nil {
		return err
	}

	if err := consumer.Unsubscribe(ctx); err != nil {
		return err
	}

	b.consumers.Del(cmd.ConsumerId)
	return nil
}

func (b *broker) Flow(ctx context.Context, cmd *geckod.CommandFlow) error {
	consumer, err := b.consumers.Get(cmd.ConsumerId)
	if err != nil {
		return nil
	}

	return consumer.Flow(ctx, cmd.MsgPermits)
}

func (b *broker) Ack(ctx context.Context, cmd *geckod.CommandAck) error {
	consumer, err := b.consumers.Get(cmd.ConsumerId)
	if err != nil {
		return err
	}

	return consumer.Ack(ctx, geckod.AckType(cmd.AckType), cmd.MessageIds)
}

func (b *broker) Send(ctx context.Context, cmd *geckod.CommandSend) error {
	producer, err := b.producers.Get(cmd.ProducerId)
	if err != nil {
		return err
	}
	seqId := producer.GetSequenceId()
	if cmd.SequenceId >= seqId {
		return errs.ErrDuplicatedSequenceId
	}

	return producer.Send(ctx, &geckod.RawMessage{
		TopicName:    cmd.TopicName,
		ProducerName: producer.Name,
		SequenceId:   cmd.SequenceId,
		Timestamp:    time.UnixMilli(cmd.Timestamp),
		Key:          cmd.Key,
		Payload:      cmd.Payload,
	})
}
