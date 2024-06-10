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
	ErrNoActorInbox                error    = errors.New("actor has no inbox")
	ErrEntryTypeNotSupported       error    = errors.New("entry type not supported")
	ErrBinaryFileNotAnImage        error    = errors.New("binary file is not an image")
	ErrBinaryFileUnsupportedFormat error    = errors.New("binary file has unsupported image format")
	ErrEntryNotFound               WebError = &webError{
		status: 404,
		msg:    "entry not found",
	}
	ErrInteractionNotFound WebError = &webError{
		status: 404,
		msg:    "interaction not found",
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
	ErrNoWebmentionEndpointFound error = errors.New("no webmention endpoint found")
	ErrTypeAlreadyRegistered     error = errors.New("type already registered")
	ErrTypeNotRegistered         error = errors.New("type not registered")
)
