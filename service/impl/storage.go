package impl

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/bitcapybara/geckod"
	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
	"go.uber.org/atomic"
)

var _ service.Storage = (*memStorage)(nil)

type memStorage struct {
	mu         sync.Mutex
	latest_id  atomic.Uint64
	msgs       map[uint64]*geckod.RawMessage
	sortedMsgs []uint64
}

// Add implements service.Storage
func (s *memStorage) Add(ctx context.Context, msg *geckod.RawMessage) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.latest_id.Inc()

	s.msgs[id] = msg
	s.sortedMsgs = append(s.sortedMsgs, id)
	sort.Slice(s.sortedMsgs, func(i, j int) bool {
		return s.sortedMsgs[i] < s.sortedMsgs[j]
	})
	return id, nil
}

// DelUntil implements service.Storage
func (s *memStorage) DelUntil(ctx context.Context, id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := s.sortedMsgs[0]
	endIdx := id - start

	for _, id := range s.sortedMsgs[:endIdx+1] {
		delete(s.msgs, id)
	}
	s.sortedMsgs = s.sortedMsgs[endIdx+1:]
	return nil
}

// Get implements service.Storage
func (s *memStorage) Get(ctx context.Context, id uint64) (*geckod.RawMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if msg, ok := s.msgs[id]; ok {
		return msg, nil
	}

	return nil, errs.ErrNotFound
}

// GetBatch implements service.Storage
func (s *memStorage) GetBatch(ctx context.Context, ids []uint64) ([]*geckod.RawMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]*geckod.RawMessage, len(ids))
	for _, id := range ids {
		if msg, ok := s.msgs[id]; ok {
			res = append(res, msg)
		}
	}

	return res, nil
}

// GetMore implements service.Storage
func (s *memStorage) GetMore(ctx context.Context, limit uint64) ([]*geckod.RawMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	endIdx := len(s.sortedMsgs)
	if int(limit) < len(s.sortedMsgs) {
		endIdx = int(limit)
	}

	ids := s.sortedMsgs[:endIdx]
	res := make([]*geckod.RawMessage, len(ids))
	for _, id := range ids {
		if msg, ok := s.msgs[id]; ok {
			res = append(res, msg)
		}
	}
	return res, nil
}

// GetRange implements service.Storage
func (s *memStorage) GetRange(ctx context.Context, from uint64, to uint64) ([]*geckod.RawMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := s.sortedMsgs[0]

	fromIdx := from - start
	if fromIdx < 0 {
		fromIdx = 0
	}
	if int(fromIdx) >= len(s.sortedMsgs) {
		return nil, errors.New("from id not found in storage")
	}

	endIdx := len(s.sortedMsgs)
	if int(to) < endIdx {
		endIdx = int(to)
	}

	res := make([]*geckod.RawMessage, endIdx-int(from))
	for _, id := range s.sortedMsgs[from : endIdx+1] {
		if msg, ok := s.msgs[id]; ok {
			res = append(res, msg)
		}
	}

	return res, nil
}
