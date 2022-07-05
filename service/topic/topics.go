package topic

import (
	"sync"

	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
)

var _ service.Topics = (*topics)(nil)

type topics struct {
	mu     sync.Mutex
	topics map[string]service.Topic
}

func (t *topics) GetOrCreate(name string) service.Topic {
	t.mu.Lock()
	defer t.mu.Unlock()

	if topic, ok := t.topics[name]; ok {
		return topic
	}

	topic := newTopic(name)
	t.topics[name] = topic

	return topic
}

func (t *topics) Get(name string) (service.Topic, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if topic, ok := t.topics[name]; ok {
		return topic, nil
	}

	return nil, errs.ErrNotFound
}

func (t *topics) Del(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.topics, name)
}
