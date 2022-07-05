package broker

import (
	"time"

	"github.com/bitcapybara/geckod"
	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
)

type Broker interface {
	Connect(*geckod.CommandConnect) (*geckod.CommandConnected, error)
	Producer(*geckod.CommandProducer) (*geckod.CommandProducerSuccess, error)
	Subscribe(*geckod.CommandSubscribe) (*geckod.CommandSubscribeSuccess, error)
	Unsubscribe(*geckod.CommandUnsubscribe) error
	Flow(*geckod.CommandFlow) error
	Ack(*geckod.CommandAck) error
}

type Authenticator = func(username, passwd string, method geckod.ConnectAuthMethod) (bool, error)

var _ Broker = (*broker)(nil)

type broker struct {
	topics    service.Topics
	clients   service.Clients
	producers service.Producers
	consumers service.Consumers
	storage   service.Storage

	authFn Authenticator
}

func (b *broker) Connect(cmd *geckod.CommandConnect) (*geckod.CommandConnected, error) {
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

func (b *broker) Producer(cmd *geckod.CommandProducer) (*geckod.CommandProducerSuccess, error) {
	topic, err := b.topics.GetOrCreate(cmd.Topic)
	if err != nil {
		return nil, err
	}

	info, err := b.producers.GetOrCreate(cmd.ClientId, cmd.ProducerName, geckod.ProducerAccessMode(cmd.AccessMode))
	if err != nil {
		return nil, err
	}

	if err := topic.AddProducer(info); err != nil {
		return nil, err
	}

	return &geckod.CommandProducerSuccess{
		ProducerId: info.Id,
	}, nil
}

func (b *broker) Subscribe(cmd *geckod.CommandSubscribe) (*geckod.CommandSubscribeSuccess, error) {
	topic, err := b.topics.GetOrCreate(cmd.Topic)
	if err != nil {
		return nil, err
	}

	consumer, err := b.consumers.Add(&service.AddConsumerParams{
		ClientId:     cmd.ClientId,
		ConsumerName: cmd.ConsumerName,
		TopicName:    cmd.Topic,
		SubName:      cmd.SubName,
		SubType:      geckod.SubScriptionType(cmd.SubType),
	})
	if err != nil {
		return nil, err
	}

	if err := topic.Subscribe(consumer); err != nil {
		return nil, err
	}

	return &geckod.CommandSubscribeSuccess{
		ConsumerId: consumer.Id,
	}, nil
}

func (b *broker) Unsubscribe(cmd *geckod.CommandUnsubscribe) error {

	consumer, err := b.consumers.Get(cmd.ConsumerId)
	if err != nil {
		return err
	}

	return consumer.Unsubscribe()
}

func (b *broker) Flow(cmd *geckod.CommandFlow) error {
	consumer, err := b.consumers.Get(cmd.ConsumerId)
	if err != nil {
		return nil
	}

	return consumer.Flow(cmd.MsgPermits)
}

func (b *broker) Ack(cmd *geckod.CommandAck) error {
	consumer, err := b.consumers.Get(cmd.ConsumerId)
	if err != nil {
		return err
	}

	return consumer.Ack(geckod.AckType(cmd.AckType), cmd.MessageIds)
}

func (b *broker) Send(cmd *geckod.CommandSend) error {
	producer, err := b.producers.Get(cmd.ProducerId)
	if err != nil {
		return err
	}
	// 存储
	if _, err := b.storage.Add(&service.RawMessage{
		TopicName:    cmd.TopicName,
		ProducerName: producer.Name,
		SequenceId:   cmd.SequenceId,
		Timestamp:    time.UnixMilli(cmd.Timestamp),
		Key:          cmd.Key,
		Payload:      cmd.Payload,
	}); err != nil {
		return err
	}

	// 更新元数据
	if err := b.producers.UpdateSequenceId(producer.Id, cmd.SequenceId); err != nil {
		return err
	}

	return producer.Send()
}
