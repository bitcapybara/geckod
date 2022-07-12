package service

import "github.com/bitcapybara/geckod"

type Storage interface {
	Add(*geckod.RawMessage) (uint64, error)
	Get(id uint64) (*geckod.RawMessage, error)
	GetRange(from uint64, to uint64) ([]*geckod.RawMessage, error)
	GetBatch(ids []uint64) ([]*geckod.RawMessage, error)
	GetMore(limit uint64) ([]*geckod.RawMessage, error)
	DelUntil(id uint64) error
}
