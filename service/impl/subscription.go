package impl

import (
	"github.com/bitcapybara/geckod"
	"github.com/bitcapybara/geckod/service"
)

var _ service.Subscription = (*subscription)(nil)

type subscription struct {
}

func newSubscription() service.Subscription {
	return nil
}
func (s *subscription) GetName() string {
	panic("not implemented") // TODO: Implement
}

func (s *subscription) GetType() geckod.SubScriptionType {
	panic("not implemented") // TODO: Implement
}

func (s *subscription) AddConsumer(_ service.Consumer) error {
	panic("not implemented") // TODO: Implement
}

func (s *subscription) DelConsumer(consumerId uint64) error {
	panic("not implemented") // TODO: Implement
}

func (s *subscription) Flow(consumerId uint64, permits uint64) error {
	panic("not implemented") // TODO: Implement
}

func (s *subscription) Ack(ackType geckod.AckType, msgIds []uint64) error {
	panic("not implemented") // TODO: Implement
}

func (s *subscription) Unsubscribe(consumerId uint64) error {
	panic("not implemented") // TODO: Implement
}

func (s *subscription) Close() error {
	panic("not implemented") // TODO: Implement
}
