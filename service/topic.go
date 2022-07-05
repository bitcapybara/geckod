package service

type Topics interface {
	GetOrCreate(name string) (Topic, error)
	Get(name string) (Topic, error)
	Del(name string) error
}

type Topic interface {
	GetName() string

	// 处理客户端生产者发送的数据
	Publish() error
	// 处理消费者订阅
	Subscribe(*Consumer) error

	// 生产者管理
	AddProducer(*Producer) error
	GetProducer(producer_id uint64) Producer
	DelProducer(producer_id uint64) error

	// 订阅
	GetSubscription(subName string) (Subscription, error)

	// 释放资源
	Close() error
}
