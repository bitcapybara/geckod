package service

import (
	"sync"

	"github.com/bitcapybara/geckod"
)

type Producers interface {
	Create(cfg *AddProducerParams) (*Producer, error)
	Get(id uint64) (*Producer, error)
	Del(id uint64)
}

type AddProducerParams struct {
	Id           uint64
	clientId     uint64
	ProducerName string
	AccessMode   geckod.ProducerAccessMode
	Topic        Topic
}

type Producer struct {
	Id   uint64
	Name string

	mu         sync.Mutex
	sequenceId uint64

	topic Topic
}

func NewProducer(id uint64, cfg *AddProducerParams) *Producer {
	return &Producer{
		Id:         cfg.Id,
		Name:       cfg.ProducerName,
		sequenceId: 0,
		topic:      cfg.Topic,
	}
}

func (p *Producer) GetTopic() Topic {
	return nil
}

func (p *Producer) Send() error {
	return p.topic.Publish()
}

func (p *Producer) SetSequenceId(seqId uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.sequenceId = seqId
}
