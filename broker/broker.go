package broker

import (
	"github.com/bitcapybara/geckod"
	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service/client"
	"github.com/bitcapybara/geckod/service/consumer"
	"github.com/bitcapybara/geckod/service/producer"
	"github.com/bitcapybara/geckod/service/subscription"
	"github.com/bitcapybara/geckod/service/topic"
)

type Broker interface {
	Connect(*geckod.CommandConnect) (geckod.CommandConnected, error)
	Producer(cmd *geckod.CommandProducer) (*geckod.CommandProducerSuccess, error)
	Subscribe(cmd *geckod.CommandSubscribe) (*geckod.CommandSubscribeSuccess, error)
	Unsubscribe(cmd *geckod.CommandUnsubscribe) error
	Flow(cmd *geckod.CommandFlow) error
}

type Authenticator = func(username, passwd string, method geckod.ConnectAuthMethod) (bool, error)

type broker struct {
	clients   client.ClientManager
	producers producer.ProducerInfoManager
	consumers consumer.ConsumerInfoManager
	topics    topic.TopicManager

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

	info, err := b.consumers.Add(&consumer.AddConsumerParams{
		ClientId:     cmd.ClientId,
		ConsumerName: cmd.ConsumerName,
		TopicName:    cmd.Topic,
		SubName:      cmd.SubName,
		SubType:      geckod.SubScriptionType(cmd.SubType),
	})
	if err != nil {
		return nil, err
	}

	if err := topic.Subscribe(info); err != nil {
		return nil, err
	}

	return &geckod.CommandSubscribeSuccess{
		ConsumerId: info.Id,
	}, nil
}

func (b *broker) Unsubscribe(cmd *geckod.CommandUnsubscribe) error {

	sub, err := b.getSubscription(cmd.ConsumerId)
	if err != nil {
		return nil
	}

	return sub.Unsubscribe(cmd.ConsumerId)
}

func (b *broker) Flow(cmd *geckod.CommandFlow) error {

	sub, err := b.getSubscription(cmd.ConsumerId)
	if err != nil {
		return nil
	}

	return sub.GetDispatcher().Flow(cmd.ConsumerId, cmd.MsgPermits)
}

func (b *broker) Ack(cmd *geckod.CommandAck) error {
	return nil
}

func (b *broker) Send(cmd *geckod.CommandSend) error {
	return nil
}

func (b *broker) getSubscription(consumerId uint64) (subscription.Subscription, error) {
	info, err := b.consumers.Get(consumerId)
	if err != nil {
		return nil, err
	}

	topic, err := b.topics.Get(info.Topic)
	if err != nil {
		return nil, err
	}

	return topic.GetSubscription(info.SubName)
}
