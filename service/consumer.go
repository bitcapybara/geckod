package service

import "github.com/bitcapybara/geckod"

type Consumers interface {
	Add(*AddConsumerParams) (*Consumer, error)
	Get(id uint64) (*Consumer, error)
	Del(id uint64) error
}

type AddConsumerParams struct {
	ClientId     uint64
	ConsumerName string
	TopicName    string
	SubName      string
	SubType      geckod.SubScriptionType
}

type Consumer struct {
	Id        uint64
	Name      string
	ClientId  string
	TopicName string

	sub Subscription
}

func (c *Consumer) Unsubscribe() error {
	return c.sub.Unsubscribe(c.Id)
}

func (c *Consumer) Flow(permits uint64) error {
	return c.sub.GetDispatcher().Flow(c.Id, permits)
}

func (c *Consumer) Ack(ackType geckod.AckType, msgIds []uint64) error {
	if err := c.sub.GetType().MatchAckType(ackType); err != nil {
		return err
	}
	return c.sub.Ack(ackType, msgIds)
}
