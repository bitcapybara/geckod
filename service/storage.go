package service

import (
	"context"

	"github.com/bitcapybara/geckod"
)

type Storage interface {
	Add(ctx context.Context, msg *geckod.RawMessage) (uint64, error)
	Get(ctx context.Context, id uint64) (*geckod.RawMessage, error)
	GetRange(ctx context.Context, from uint64, to uint64) ([]*geckod.RawMessage, error)
	GetBatch(ctx context.Context, ids []uint64) ([]*geckod.RawMessage, error)
	GetMore(ctx context.Context, limit uint64) ([]*geckod.RawMessage, error)
	DelUntil(ctx context.Context, id uint64) error
}
