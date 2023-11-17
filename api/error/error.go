package error

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type Error struct {
	Status
	cause error
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

func (e *Error) Error() string {
	return e.Message
}

func New(code int, reason, message string) *Error {
	return &Error{
		Status: Status{
			Code:    int32(code),
			Message: message,
			Reason:  reason,
		},
	}
}

var AuthError = New(int(TypeStatusCode_AUTH_ERROR.Number()), "", "")
var AuthAppIdError = New(int(TypeStatusCode_AUTH_APPID_ERROR.Number()), "", "")

func (e *Error) Marshal() (res []byte) {
	res, _ = proto.Marshal(e)
	return res
}
