package app

import "errors"

type WebError interface {
	Error() string
	Status() int
}

type webError struct {
	status int
	msg    string
}

func (e *webError) Error() string {
	return e.msg
}

func (e *webError) Status() int {
	return e.status
}

var (
	ErrNoActorInbox          error    = errors.New("actor has no inbox")
	ErrEntryTypeNotSupported error    = errors.New("entry type not supported")
	ErrEntryNotFound         WebError = &webError{
		status: 404,
		msg:    "entry not found",
	}
	ErrUnsupportedObjectType WebError = &webError{
		status: 501,
		msg:    "object type not supported",
	}
	ErrUnsupportedActionType WebError = &webError{
		status: 501,
		msg:    "action type not supported",
	}
	ErrConflictingId WebError = &webError{
		status: 409,
		msg:    "object with same id but different type already exists",
	}
	ErrUnableToVerifySignature WebError = &webError{
		status: 403,
		msg:    "cannot verify signature",
	}
)
