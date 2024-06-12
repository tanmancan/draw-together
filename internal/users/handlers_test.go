package users

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/auth"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mockUserRepository() {
	MakeUserRepository(MockRepository)
}

func resetMock() {
	MakeUserRepository(nil)
}

func TestHandleGetUser(t *testing.T) {
	mockUserRepository()
	defer resetMock()

	reqID := uuid.New()
	user := &model.User{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		Name:      uuid.NewString(),
		CreatedAt: timestamppb.Now(),
	}
	userInvalidID := &model.User{
		Id: &model.UUID{
			Value: "invalid-uuid",
		},
		Name:      uuid.NewString(),
		CreatedAt: timestamppb.Now(),
	}
	ctxWithReq := context.WithValue(context.Background(), helpers.RequestID, reqID.String())
	ctxWithUser := helpers.AddContextUser(ctxWithReq, user)
	ctxWithInvalidUserId := helpers.AddContextUser(ctxWithReq, userInvalidID)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		want     *model.User
		wantErr  bool
		setUp    func()
		tearDown func()
	}{
		{
			name: "HandleGetUser returns an error if context does not contain user",
			args: args{
				ctx: ctxWithReq,
			},
			wantErr: true,
		},
		{
			name: "HandleGetUser returns an error if user contains an invalid UUID",
			args: args{
				ctx: ctxWithInvalidUserId,
			},
			wantErr: true,
		},
		{
			name: "HandleGetUser returns an error if user UUID can't be parsed",
			args: args{
				ctx: ctxWithUser,
			},
			wantErr: true,
		},
		{
			name: "HandleGetUser returns an error if Repository.GetByID returns an error",
			args: args{
				ctx: ctxWithUser,
			},
			wantErr: true,
			setUp: func() {
				MockRepository.ErrGetById = errors.New("test error")
			},
			tearDown: func() {
				MockRepository.Reset()
			},
		},
		{
			name: "HandleGetUser returns an user if validation passes",
			args: args{
				ctx: ctxWithUser,
			},
			want: user,
			setUp: func() {
				MockRepository.RetGetById = user
			},
			tearDown: func() {
				MockRepository.Reset()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp()
			}
			defer func() {
				if tt.tearDown != nil {
					tt.tearDown()
				}
			}()
			got, err := HandleGetUser(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleGetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleGetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateHasCtxUser(t *testing.T) {
	u := &model.User{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		Name:      uuid.NewString(),
		CreatedAt: timestamppb.Now(),
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "validateHasCtxUser returns the user if it exists in context",
			args: args{
				ctx: helpers.AddContextUser(context.Background(), u),
			},
			want: u,
		},
		{
			name: "validateHasCtxUser returns an error if the user does not exists in context",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateHasCtxUser(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHasCtxUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateHasCtxUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateHasName(t *testing.T) {
	testName := uuid.NewString()
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "validateHasName returns name if it exists",
			args: args{
				ctx:  context.Background(),
				name: testName,
			},
			want: testName,
		},
		{
			name: "validateHasName returns error if name is empty",
			args: args{
				ctx:  context.Background(),
				name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateHasName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHasName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateHasName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleCreateUser(t *testing.T) {
	ctx := context.Background()
	testUserName := uuid.NewString()
	u := &model.User{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		Name:      testUserName,
		CreatedAt: timestamppb.Now(),
	}
	type args struct {
		ctx context.Context
		r   *model.CreateUserRequest
	}
	tests := []struct {
		name     string
		args     args
		want     string
		want1    string
		wantErr  bool
		setUp    func()
		tearDown func()
	}{
		{
			name: "HandleCreateUser returns a user and token if successful",
			args: args{
				ctx: ctx,
				r: &model.CreateUserRequest{
					Name: testUserName,
				},
			},
			want: testUserName,
			setUp: func() {
				MakeUserRepository(MockRepository)
				MockRepository.RetCreateUser = u
			},
			tearDown: func() {
				MockRepository.Reset()
			},
		},
		{
			name: "HandleCreateUser returns an error if user already exists in context",
			args: args{
				ctx: helpers.AddContextUser(context.Background(), u),
				r: &model.CreateUserRequest{
					Name: testUserName,
				},
			},
			wantErr: true,
		},
		{
			name: "HandleCreateUser returns an error if request parameter is invalid",
			args: args{
				ctx: context.Background(),
				r: &model.CreateUserRequest{
					Name: "",
				},
			},
			wantErr: true,
		},
		{
			name: "HandleCreateUser returns an error if Repository.CreateUser returns an error",
			args: args{
				ctx: context.Background(),
				r: &model.CreateUserRequest{
					Name: testUserName,
				},
			},
			wantErr: true,
			setUp: func() {
				MakeUserRepository(MockRepository)
				MockRepository.ErrCreateUser = status.Error(codes.Internal, "this is expected error")
			},
			tearDown: func() {
				MockRepository.Reset()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp()
			}
			defer func() {
				if tt.tearDown != nil {
					tt.tearDown()
				}
			}()
			got, got1, err := HandleCreateUser(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}

			if got.Name != tt.want {
				t.Errorf("HandleCreateUser() got.Name = %v, want %v", got, tt.want)
			}

			uid, err := auth.JwtAuthHandler(ctx, got1)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.JwtAuthHandler(ctx, got1) error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotID, _ := helpers.ProtoToUUID(got.Id)
			if gotID != uid {
				t.Errorf("HandleCreateUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHandleDeleteUser(t *testing.T) {
	userID := uuid.New()
	u := &model.User{
		Id: &model.UUID{
			Value: userID.String(),
		},
		Name:      uuid.NewString(),
		CreatedAt: timestamppb.Now(),
	}
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		setUp    func()
		tearDown func()
	}{
		{
			name: "HandleDeleteUser returns nil if no errors occur",
			args: args{
				ctx:    helpers.AddContextUser(context.Background(), u),
				userID: userID,
			},
			setUp: func() {
				MakeUserRepository(MockRepository)
				MockRepository.RetGetById = u
			},
			tearDown: func() {
				MockRepository.Reset()
			},
		},
		{
			name: "HandleDeleteUser returns error if user does not exist in store",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			wantErr: true,
		},
		{
			name: "HandleDeleteUser returns error if Repository.DeleteUser returns an error",
			args: args{
				ctx:    helpers.AddContextUser(context.Background(), u),
				userID: userID,
			},
			wantErr: true,
			setUp: func() {
				MakeUserRepository(MockRepository)
				MockRepository.RetGetById = u
				MockRepository.ErrDeleteUser = errors.New("test error")
			},
			tearDown: func() {
				MockRepository.Reset()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp()
			}
			defer func() {
				if tt.tearDown != nil {
					tt.tearDown()
				}
			}()
			if err := HandleDeleteUser(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("HandleDeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
