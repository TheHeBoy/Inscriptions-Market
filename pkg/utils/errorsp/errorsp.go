package errorsp

import "github.com/pkg/errors"

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func WithStack(err error) error {
	_, ok := err.(stackTracer)
	if ok {
		return err
	}

	return errors.WithStack(err)
}
