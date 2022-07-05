package service

import (
	"time"
)

type Storage interface {
	Add(*RawMessage) (uint64, error)
	Get(id uint64) (*RawMessage, error)
	GetRange(from uint64, to uint64) ([]*RawMessage, error)
	GetBatch(ids []uint64) ([]*RawMessage, error)
	DelUntil(id uint64) error
}

type RawMessage struct {
	Id           uint64
	TopicName    string
	ProducerName string
	SequenceId   uint64
	Timestamp    time.Time
	Key          string
	Payload      []byte
}
