package helpers

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create a new error based on "status.Error" provided by the grpc library
// Wrapping all errors around the grpc error allows us to take advantage of
// standard error codes and messages.
func MakeError(ctx context.Context, c codes.Code, msg string, args ...any) error {
	reqID := GetReqIdFromContext(ctx)

	switch c {
	case codes.Canceled:
		msg = fmt.Sprintf("%s: %s", "canceled", msg)
	case codes.Unknown:
		msg = fmt.Sprintf("%s: %s", "unknown", msg)
	case codes.InvalidArgument:
		msg = fmt.Sprintf("%s: %s", "invalid argument", msg)
	case codes.DeadlineExceeded:
		msg = fmt.Sprintf("%s: %s", "deadline exceeded", msg)
	case codes.NotFound:
		msg = fmt.Sprintf("%s: %s", "not found", msg)
	case codes.AlreadyExists:
		msg = fmt.Sprintf("%s: %s", "already exists", msg)
	case codes.PermissionDenied:
		msg = fmt.Sprintf("%s: %s", "permission denied", msg)
	case codes.ResourceExhausted:
		msg = fmt.Sprintf("%s: %s", "resource exhausted", msg)
	case codes.FailedPrecondition:
		msg = fmt.Sprintf("%s: %s", "failed precondition", msg)
	case codes.Aborted:
		msg = fmt.Sprintf("%s: %s", "aborted", msg)
	case codes.OutOfRange:
		msg = fmt.Sprintf("%s: %s", "out of range", msg)
	case codes.Unimplemented:
		msg = fmt.Sprintf("%s: %s", "unimplemented", msg)
	case codes.Internal:
		msg = fmt.Sprintf("%s: %s", "internal", msg)
	case codes.Unavailable:
		msg = fmt.Sprintf("%s: %s", "unavailable", msg)
	case codes.DataLoss:
		msg = fmt.Sprintf("%s: %s", "data loss", msg)
	case codes.Unauthenticated:
		msg = fmt.Sprintf("%s: %s", "unauthenticated", msg)
	}
	msg = fmt.Sprintf("%s request_id: %s", msg, reqID)
	return status.Errorf(c, msg, args...)
}
