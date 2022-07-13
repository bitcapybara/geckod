package impl

import (
	"context"
	"sync"

	"github.com/bitcapybara/geckod"
	"github.com/bitcapybara/geckod/service"
	"github.com/bits-and-blooms/bitset"
	"go.uber.org/atomic"
)

var _ service.Cursor = (*cursor)(nil)

// 记录消息的消费情况
type cursor struct {
	storage service.Storage

	mu sync.Mutex
	// 当前已经发送的最大消息id
	readPosition atomic.Uint64
	// 已ack的消息，标记为1
	bits bitset.BitSet
}

// Delete implements service.Cursor
func (c *cursor) Delete(ctx context.Context, msgIds []uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, id := range msgIds {
		c.bits.Set(uint(id))
	}

	return nil
}

// GetMoreMessages implements service.Cursor
func (c *cursor) GetMoreMessages(ctx context.Context, num uint64) ([]*geckod.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	msgs, err := c.storage.GetMore(ctx, num)
	if err != nil {
		return nil, err
	}
	c.readPosition.Add(uint64(len(msgs)))

	return msgs, nil
}

// MarkDelete implements service.Cursor
func (c *cursor) MarkDelete(ctx context.Context, msgId uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bits.Set(uint(msgId))
	for i, e := c.bits.NextClear(0); e && i < uint(msgId); i, e = c.bits.NextClear(i + 1) {
		c.bits.Set(i)
	}
	return nil
}
