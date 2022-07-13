package impl

import (
	"context"
	"errors"
	"sync"

	"github.com/bitcapybara/geckod"
	cmdpb "github.com/bitcapybara/geckod-proto/gen/go/proto/command"
	"github.com/bitcapybara/geckod/errs"
	"github.com/bitcapybara/geckod/service"
)

var _ service.Subscription = (*subscription)(nil)

type subscription struct {
	name    string
	subType geckod.SubScriptionType

	topic  service.Topic
	cursor service.Cursor

	mu         sync.Mutex
	dispatcher service.Dispatcher
}

func newSubscription() service.Subscription {
	return nil
}

func (s *subscription) GetName() string {
	return s.name
}

func (s *subscription) GetType() geckod.SubScriptionType {
	return s.subType
}

func (s *subscription) AddConsumer(ctx context.Context, consumer *service.Consumer) error {
	// 如果没有 dispatcher，则创建一个新的
	return s.getOrCreateDispatcher().AddConsumer(ctx, consumer)
}

func (s *subscription) getOrCreateDispatcher() service.Dispatcher {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.dispatcher == nil {
		s.dispatcher = newDispatcher()
	}

	return s.dispatcher
}

func (s *subscription) DelConsumer(ctx context.Context, consumerId uint64) error {
	//
	dispatcher, err := s.GetDispatcher()
	if err != nil {
		if err != errs.ErrNotFound {
			return err
		}
		return nil
	}

	if err := dispatcher.DelConsumer(ctx, consumerId); err != nil {
		return err
	}

	if len(dispatcher.GetConsumers(ctx)) != 0 {
		return nil
	}

	// 清理资源
	if err := dispatcher.Close(ctx); err != nil {
		return err
	}

	return s.topic.RemoveSubscription(ctx, s.name)
}

func (s *subscription) Flow(ctx context.Context, consumerId uint64, permits uint64) error {
	dispatcher, err := s.GetDispatcher()
	if err != nil {
		return err
	}
	return dispatcher.Flow(ctx, consumerId, permits)
}

func (s *subscription) Ack(ctx context.Context, ackType geckod.AckType, msgIds []uint64) error {
	if ackType == geckod.AckType(cmdpb.Ack_Cumulative) {
		if len(msgIds) != 1 {
			return errors.New("Invalid cumulative ack received with multiple message ids")
		}
		if err := s.cursor.MarkDelete(ctx, msgIds[0]); err != nil {
			return err
		}
		return nil
	}

	return s.cursor.Delete(ctx, msgIds)
}

func (s *subscription) Unsubscribe(ctx context.Context, consumer *service.Consumer) error {
	dispatcher, err := s.GetDispatcher()
	if err != nil {
		return nil
	}
	if !dispatcher.CanUnsubscribe(ctx, consumer.Id) {
		return errors.New("Shared consumer attempting to unsubscribe")
	}

	if err := consumer.Close(ctx); err != nil {
		return err
	}

	return s.Close(ctx)
}

func (s *subscription) GetDispatcher() (service.Dispatcher, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.dispatcher != nil {
		return s.dispatcher, nil
	}
	return nil, errs.ErrNotFound
}

func (s *subscription) Close(ctx context.Context) error {
	if err := s.topic.Unsubscribe(ctx, s.name); err != nil {
		return err
	}

	dispatcher, err := s.GetDispatcher()
	if err != nil {
		return err
	}
	if err := dispatcher.Close(ctx); err != nil {
		return err
	}
	return nil
}
