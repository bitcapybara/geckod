package errs

import (
	"errors"

	cmdpb "github.com/bitcapybara/geckod-proto/gen/go/proto/command"
)

type Error cmdpb.Error

func (c *Error) Error() string {
	return c.Message
}

func NewCommandError(code cmdpb.Error_Code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

var (
	ErrNotFound = errors.New("not found")

	ErrAuthFailed                = NewCommandError(cmdpb.Error_AuthenticateFail, "身份校验失败")
	ErrUnmatchAckType            = NewCommandError(cmdpb.Error_UnmatchAckType, "cumulative ack on shared subscription")
	ErrProducerAlreadyExists     = NewCommandError(cmdpb.Error_ProducerAlreadyExists, "生产者已存在")
	ErrConsumerAlreadyExists     = NewCommandError(cmdpb.Error_ConsumerAlreadyExists, "消费者已存在")
	ErrProducerExclusiveConflict = errors.New("")
)
