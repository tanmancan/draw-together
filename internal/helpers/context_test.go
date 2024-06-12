package helpers

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetRequestIdFromCtx(t *testing.T) {
	reqUuid := uuid.New()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetRequestIdFromCtx should return a valid UUID string if context has reqID",
			args: args{
				ctx: context.WithValue(context.Background(), RequestID, reqUuid.String()),
			},
			want: reqUuid.String(),
		},
		{
			name: "GetRequestIdFromCtx should return a not set value if context does not have reqID",
			args: args{
				ctx: context.Background(),
			},
			want: reqIdNotSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetReqIdFromContext(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestIdFromCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateContextWithRequestID(t *testing.T) {
	t.Run("CreateContextWithRequestID should always return a context with reqID",
		func(t *testing.T) {
			ctx := context.Background()
			got := CreateContextWithRequestID(ctx)
			gotReqID := got.Value(RequestID)
			if gotReqID == "" || gotReqID == nil {
				t.Errorf("CreateContextWithRequestID() = %v, want %v", got, "[UUID]")
			}
		})

	t.Run("CreateContextWithRequestID should return a context with new reqID if one already exist",
		func(t *testing.T) {
			oldReqID := uuid.NewString()
			ctxWithReqID := context.WithValue(context.Background(), RequestID, oldReqID)
			got := CreateContextWithRequestID(ctxWithReqID)
			gotReqID := got.Value(RequestID)
			if gotReqID == "" || gotReqID == nil {
				t.Errorf("CreateContextWithRequestID() = %v, oldReqId = %v, want %v", got, oldReqID, "[UUID]")
			}
			if gotReqID == oldReqID {
				t.Errorf("CreateContextWithRequestID() = %v, oldReqId = %v, want %v", got, oldReqID, "[New UUID]")
			}
		})
}

func TestGetUserFromContext(t *testing.T) {
	ctxWithUser := context.Background()
	testUser := model.User{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		Name:      uuid.NewString(),
		CreatedAt: timestamppb.Now(),
	}
	ctxWithUser = context.WithValue(ctxWithUser, AuthUserID, testUser.Id.Value)
	ctxWithUser = context.WithValue(ctxWithUser, AuthUserName, testUser.Name)
	ctxWithUser = context.WithValue(ctxWithUser, AuthUserCreatedAt, testUser.CreatedAt.String())
	ctxWithoutUser := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *model.User
	}{
		{
			name: "GetUserFromContext returns a model.User if it exists",
			args: args{
				ctx: ctxWithUser,
			},
			want: &testUser,
		},
		{
			name: "GetUserFromContext returns nil if user does not exist in context",
			args: args{
				ctx: ctxWithoutUser,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUserFromContext(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReqIdFromContext(t *testing.T) {
	reqID := uuid.New()
	ctxWithReqId := context.WithValue(context.Background(), RequestID, reqID.String())
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetReqIdFromContext will return reqID if it exists",
			args: args{
				ctx: ctxWithReqId,
			},
			want: reqID.String(),
		},
		{
			name: "GetReqIdFromContext will return default string if it does not exist",
			args: args{
				ctx: context.Background(),
			},
			want: reqIdNotSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetReqIdFromContext(tt.args.ctx); got != tt.want {
				t.Errorf("GetReqIdFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddContextUser(t *testing.T) {
	ctx := context.Background()
	testUser := &model.User{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		Name:      uuid.NewString(),
		CreatedAt: timestamppb.Now(),
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "AddContextUser adds user fields to context",
			args: args{
				ctx:  ctx,
				user: testUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddContextUser(tt.args.ctx, tt.args.user)
			u := GetUserFromContext(got)
			if !reflect.DeepEqual(testUser, u) {
				t.Errorf("AddContextUser() = %v, want %v", testUser, u)
			}
		})
	}
}
