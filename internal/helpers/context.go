package helpers

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ctxKey int

const (
	RequestID ctxKey = iota
	AuthUserID
	AuthUserName
	AuthUserCreatedAt
)

var reqIdNotSet = "req_id_not_set"

func CreateContextWithRequestID(ctx context.Context) context.Context {
	reqID := uuid.NewString()
	return context.WithValue(ctx, RequestID, reqID)
}

func GetReqIdFromContext(ctx context.Context) string {
	reqID := ctx.Value(RequestID)
	if reqID == nil {
		return reqIdNotSet
	}
	return reqID.(string)
}

func GetUserFromContext(ctx context.Context) *model.User {
	var ts timestamppb.Timestamp
	id := ctx.Value(AuthUserID)
	name := ctx.Value(AuthUserName)
	createdAt := ctx.Value(AuthUserCreatedAt)
	if id != nil && name != nil && createdAt != nil {
		prototext.Unmarshal([]byte(createdAt.(string)), &ts)
		return &model.User{
			Id: &model.UUID{
				Value: id.(string),
			},
			Name:      name.(string),
			CreatedAt: &ts,
		}
	}

	return nil
}

func AddContextUser(ctx context.Context, user *model.User) context.Context {
	ctx = context.WithValue(ctx, AuthUserID, user.Id.Value)
	ctx = context.WithValue(ctx, AuthUserName, user.Name)
	ctx = context.WithValue(ctx, AuthUserCreatedAt, user.CreatedAt.String())
	return ctx
}
