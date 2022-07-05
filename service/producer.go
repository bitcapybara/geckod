package service

import (
	"github.com/bitcapybara/geckod"
)

type Producers interface {
	GetOrCreate(client_id uint64, name string, access_mode geckod.ProducerAccessMode) (*Producer, error)
	Get(id uint64) (*Producer, error)
	Del(id uint64) (*Producer, error)
	UpdateSequenceId(id, seqId uint64) error
}

type Producer struct {
	Id         uint64
	Name       string
	SequenceId uint64

	topic Topic
}

func (p *Producer) GetTopic() Topic {
	return nil
}

func (p *Producer) Send() error {
	return p.topic.Publish()
}
