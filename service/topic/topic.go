package topic

import "github.com/bitcapybara/geckod/service"

var _ service.Topic = (*topic)(nil)

type topic struct {
}

func newTopic(name string) *topic {
	return nil
}

func (t *topic) GetName() string {
	panic("not implemented") // TODO: Implement
}

// 处理客户端生产者发送的数据
func (t *topic) Publish() error {
	panic("not implemented") // TODO: Implement
}

// 处理消费者订阅
func (t *topic) Subscribe() (*service.Consumer, error) {
	panic("not implemented") // TODO: Implement
}

// 生产者管理
func (t *topic) AddProducer(_ *service.Producer) error {
	panic("not implemented") // TODO: Implement
}

func (t *topic) GetProducer(producer_id uint64) service.Producer {
	panic("not implemented") // TODO: Implement
}

func (t *topic) DelProducer(producer_id uint64) error {
	panic("not implemented") // TODO: Implement
}

// 订阅
func (t *topic) GetSubscription(subName string) (service.Subscription, error) {
	panic("not implemented") // TODO: Implement
}

// 释放资源
func (t *topic) Close() error {
	panic("not implemented") // TODO: Implement
}
