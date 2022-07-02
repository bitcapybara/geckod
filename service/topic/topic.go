package topic

import (
	"github.com/bitcapybara/geckod/service/consumer"
	"github.com/bitcapybara/geckod/service/producer"
	"github.com/bitcapybara/geckod/service/subscription"
)

type TopicManager interface {
	GetOrCreate(name string) (Topic, error)
	Get(name string) (Topic, error)
	Del(name string) error
}

type Topic interface {
	GetName() string

	// 处理客户端生产者发送的数据
	Publish() error
	// 处理消费者订阅
	Subscribe(*consumer.ConsumerInfo) error

	// 生产者管理
	AddProducer(*producer.ProducerInfo) error
	GetProducer(producer_id uint64) producer.Producer
	DelProducer(producer_id uint64) error

	// 订阅
	GetSubscription(subName string) (subscription.Subscription, error)

	// 释放资源
	Close() error
}
