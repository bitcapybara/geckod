package producer

import (
	"sync"

	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
	"go.uber.org/atomic"
)

var _ service.Producers = (*producers)(nil)

type producers struct {
	latestId atomic.Uint64 // init with 1

	mu              sync.Mutex
	producers       map[uint64]*service.Producer
	producersByName map[string]struct{}
}

func (p *producers) Create(cfg *service.ProducerConfig) (*service.Producer, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.producersByName[cfg.ProducerName]; ok {
		return nil, errs.ErrProducerAlreadyExists
	}

	id := p.latestId.Inc()
	producer := service.NewProducer(id, cfg)
	p.producers[id] = producer
	p.producersByName[cfg.ProducerName] = struct{}{}

	return producer, nil
}

func (p *producers) Get(id uint64) (*service.Producer, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p, ok := p.producers[id]; ok {
		return p, nil
	}

	return nil, errs.ErrNotFound
}

func (p *producers) Del(id uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if producer, ok := p.producers[id]; ok {
		delete(p.producersByName, producer.Name)
	}

	delete(p.producers, id)
}
