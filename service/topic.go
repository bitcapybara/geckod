package service

type Topics interface {
	GetOrCreate(name string) Topic
	Get(name string) (Topic, error)
	Del(name string)
}

type Topic interface {
	GetName() string

	// 处理客户端生产者发送的数据
	Publish() error
	// 处理消费者订阅
	// 生成 consumer，添加到 subscription，返回 consumer
	Subscribe(*SubscriptionOption) (*Consumer, error)

	// 生产者管理
	AddProducer(*Producer) error
	GetProducer(producer_id uint64) (*Producer, error)
	DelProducer(producer_id uint64) error

	// 订阅
	GetSubscription(subName string) (Subscription, error)

	// 释放资源
	Close() error
}

type SubscriptionOption struct {
	SubName      string
	ClientId     uint64
	ConsumerName string
	TopicName    string
}
