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

	ErrAuthFailed = NewCommandError(cmdpb.Error_AuthenticateFail, "身份校验失败")
)
