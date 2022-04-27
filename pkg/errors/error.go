package errors

var ErrMsgQueue = NewCometError("user message queue error")

func NewCometError(message string) *CometError {
	return &CometError{
		Message: message,
	}
}

type CometError struct {
	Message string
}

func (*CometError) Error() string {
	return "CometError"
}
