package impl

import (
	"sync"

	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
	"go.uber.org/atomic"
)

var _ service.Dispatcher = (*dispatcher)(nil)

type dispatcher struct {
	mu        sync.Mutex
	consumers map[uint64]*service.Consumer

	totalAvailablePermits atomic.Uint64

	sendMu sync.Mutex
	cursor service.Cursor
}

func newDispatcher() *dispatcher {
	return nil
}

func (d *dispatcher) AddConsumer(consumer *service.Consumer) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.consumers[consumer.Id]; ok {
		return errs.ErrConsumerAlreadyExists
	}

	d.consumers[consumer.Id] = consumer
	return nil
}

func (d *dispatcher) DelConsumer(consumerId uint64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.consumers, consumerId)
	return nil
}

func (d *dispatcher) GetConsumers() []*service.Consumer {
	d.mu.Lock()
	defer d.mu.Unlock()

	res := make([]*service.Consumer, 0, len(d.consumers))
	for _, consumer := range d.consumers {
		res = append(res, consumer)
	}

	return res
}

func (d *dispatcher) Flow(consumerId uint64, permits uint64) error {
	if !d.hasConsumer(consumerId) {
		return errs.ErrNotFound
	}

	d.totalAvailablePermits.Store(d.totalAvailablePermits.Load() + permits)

	return d.sendMessagesToConsumers()
}

func (d *dispatcher) sendMessagesToConsumers() error {
	// 如果当前正在发送，则直接退出
	if !d.sendMu.TryLock() {
		return nil
	}
	defer d.sendMu.Unlock()

	for {
		// 是否有可发送的消费者
		firstAvailablePermits := d.getFirstAvailableConsumerPermits()
		if firstAvailablePermits <= 0 {
			return nil
		}

		// 获取下一个可用的消费者
		consumer := d.getNextConsumer()
		if consumer == nil {
			return nil
		}

		// 要发送给消费者的消息数量
		permits := consumer.Permits()
		// 获取要发送的消息
		msgs, err := d.cursor.GetMoreMessages(permits)
		if err != nil {
			return err
		}
		// 发送
		if err := consumer.SendMessages(msgs); err != nil {
			return err
		}

		d.totalAvailablePermits.Add(-permits)
	}
}

// 获取下一个可用的消费者（优先级）
// 没有返回 nil
func (d *dispatcher) getNextConsumer() *service.Consumer {
	return nil
}

func (d *dispatcher) getFirstAvailableConsumerPermits() uint64 {
	return 0
}

func (d *dispatcher) hasConsumer(consumerId uint64) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.consumers[consumerId]; ok {
		return true
	}
	return false
}

func (d *dispatcher) CanUnsubscribe(consumerId uint64) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.consumers) == 1 && d.consumers[consumerId] != nil
}

func (d *dispatcher) SendMessages() error {
	panic("not implemented") // TODO: Implement
}

func (d *dispatcher) Close() error {
	panic("not implemented") // TODO: Implement
}
