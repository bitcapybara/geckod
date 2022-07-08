package impl

import "github.com/bitcapybara/geckod/service"

var _ service.Dispatcher = (*dispatcher)(nil)

type dispatcher struct {
}

func newDispatcher() *dispatcher {
	return nil
}

func (d *dispatcher) AddConsumer(_ *service.Consumer) error {
	panic("not implemented") // TODO: Implement
}

func (d *dispatcher) DelConsumer(consumerId uint64) error {
	panic("not implemented") // TODO: Implement
}

func (d *dispatcher) GetConsumers() []*service.Consumer {
	return nil
}

func (d *dispatcher) Flow(consumerId uint64, permits uint64) error {
	panic("not implemented") // TODO: Implement
}

func (d *dispatcher) SendMessages() error {
	panic("not implemented") // TODO: Implement
}

func (d *dispatcher) Close() error {
	panic("not implemented") // TODO: Implement
}
