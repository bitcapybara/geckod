package impl

import (
	"sync"

	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
	"go.uber.org/atomic"
)

var _ service.Consumers = (*consumers)(nil)

type consumers struct {
	latestId atomic.Uint64

	mu              sync.Mutex
	consumers       map[uint64]*service.Consumer
	consumersByName map[string]struct{}
}

func (c *consumers) GetOrCreate(params *service.AddConsumerParams) (*service.Consumer, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.consumersByName[params.ConsumerName]; ok {
		return nil, errs.ErrConsumerAlreadyExists
	}

	id := c.latestId.Inc()
	consumer := service.NewConsumer(id, params)
	c.consumers[id] = consumer
	c.consumersByName[params.ConsumerName] = struct{}{}

	return consumer, nil
}

func (c *consumers) Get(id uint64) (*service.Consumer, error) {
	panic("not implemented") // TODO: Implement
}

func (c *consumers) Del(id uint64) {
	panic("not implemented") // TODO: Implement
}

func (c *consumers) Add(consumer *service.Consumer) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.consumersByName[consumer.Name]; ok {
		return errs.ErrConsumerAlreadyExists
	}

	id := c.latestId.Inc()
	c.consumers[id] = consumer
	c.consumersByName[consumer.Name] = struct{}{}

	return nil
}
