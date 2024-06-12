package middleware

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/auth"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockServerTransportStream struct {
	md *metadata.MD
}

func getHeader(s mockServerTransportStream) (string, error) {
	if s.md == nil {
		return "", errors.New("metadata not set")
	}

	csp := s.md.Get(helpers.HeaderContentSecurityPolicy)

	if len(csp) == 0 {
		return "", errors.New("header not found")
	}

	return csp[0], nil
}

func (s *mockServerTransportStream) Method() string {
	return ""
}
func (s *mockServerTransportStream) SetHeader(md metadata.MD) error {
	s.md = &md
	return nil
}
func (s *mockServerTransportStream) SendHeader(md metadata.MD) error {
	return nil
}
func (s *mockServerTransportStream) SetTrailer(md metadata.MD) error {
	return nil
}

var ctx = context.WithValue(context.Background(), helpers.RequestID, uuid.NewString())

func TestGetChainUnaryInterceptor(t *testing.T) {
	var called = 0
	var errorCt = 0
	mockServer := func(o any) {
		switch o.(type) {
		case grpc.ServerOption:
			called++
		default:
			errorCt++
		}
	}
	tests := []struct {
		name string
	}{
		{
			name: "GetChainUnaryInterceptor returns a grpc.ServerOption",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetChainUnaryInterceptor()
			mockServer(got)
			if called == 0 {
				t.Errorf("GetChainUnaryInterceptor() = %v", got)
			}
		})
	}
}

func Test_interceptorCallsHandler(t *testing.T) {

	testInfo := &grpc.UnaryServerInfo{
		FullMethod: "testMethod",
	}
	wantResp := "test_handler_response"
	testHandler := func(context.Context, any) (any, error) {
		return wantResp, nil
	}
	type args struct {
		ctx      context.Context
		req      any
		info     *grpc.UnaryServerInfo
		handler  grpc.UnaryHandler
		unaryInt grpc.UnaryServerInterceptor
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name: "unprotectedInterceptor calls the handler",
			args: args{
				ctx:      ctx,
				req:      nil,
				info:     testInfo,
				handler:  testHandler,
				unaryInt: unprotectedInterceptor,
			},
			wantResp: wantResp,
			wantErr:  false,
		},
		{
			name: "requestIdInterceptor calls the handler",
			args: args{
				ctx:      ctx,
				req:      nil,
				info:     testInfo,
				handler:  testHandler,
				unaryInt: requestIdInterceptor,
			},
			wantResp: wantResp,
			wantErr:  false,
		},
		{
			name: "loggerInterceptor calls the handler",
			args: args{
				ctx:      ctx,
				req:      nil,
				info:     testInfo,
				handler:  testHandler,
				unaryInt: loggerInterceptor,
			},
			wantResp: wantResp,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.unaryInt(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("%v() error = %v, wantErr %v", tt.args.unaryInt, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("%v() = %v, want %v", tt.args.unaryInt, gotResp, tt.wantResp)
			}
		})
	}
}

func Test_authContextInterceptor(t *testing.T) {
	testInfoValidMethod := grpc.UnaryServerInfo{
		FullMethod: "testMethod",
	}
	testInfoInvalidMethod := grpc.UnaryServerInfo{
		FullMethod: "invalidTestMethod",
	}
	testResp := "testResp"
	testInt := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return testResp, nil
	}
	testProtectedMethods := protectedMethodsInterceptors{
		"testMethod": testInt,
	}
	type args struct {
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name: "authContextInterceptor calls wrapperInt if valid method is provided",
			args: args{
				ctx:     ctx,
				req:     nil,
				info:    &testInfoValidMethod,
				handler: nil,
			},
			wantResp: testResp,
			wantErr:  false,
		},
		{
			name: "authContextInterceptor returns an error if valid method is not found",
			args: args{
				ctx:     ctx,
				req:     nil,
				info:    &testInfoInvalidMethod,
				handler: nil,
			},
			wantResp: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpMethod := protectedMethods
			protectedMethods = testProtectedMethods
			gotResp, err := authContextInterceptor(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			protectedMethods = tmpMethod
			if (err != nil) != tt.wantErr {
				t.Errorf("authContextInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("authContextInterceptor() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_jwtProtectedInterceptor(t *testing.T) {
	ctxToken := context.Background()
	testUserID := uuid.NewString()
	testUserName := uuid.NewString()
	testCreatedAt := timestamppb.Now()
	testValidToken, _ := auth.GenerateUserToken(ctxToken, testUserID, testUserName)
	testUser := &model.User{
		Id: &model.UUID{
			Value: testUserID,
		},
		Name:      testUserName,
		CreatedAt: testCreatedAt,
	}

	mdMap := make(map[string]string)
	mdInvalidToken := make(map[string]string)
	mdInvalidToken[helpers.HeaderXRequestToken] = "testInvalidToken"
	mdValidToken := make(map[string]string)
	mdValidToken[helpers.HeaderXRequestToken] = testValidToken

	ctxNoMd := context.Background()
	ctxMdNoToken := metadata.NewIncomingContext(context.Background(), metadata.New(mdMap))
	ctxMdInvalidToken := metadata.NewIncomingContext(context.Background(), metadata.New(mdInvalidToken))
	ctxMdValidToken := metadata.NewIncomingContext(ctxToken, metadata.New(mdValidToken))

	testHandler := func(ctx context.Context, req any) (any, error) {
		return helpers.GetUserFromContext(ctx), nil
	}
	type args struct {
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
		before   func()
		after    func()
	}{
		{
			name: "jwtProtectedInterceptor returns an error if no metadata provided in context",
			args: args{
				ctx:     ctxNoMd,
				req:     nil,
				info:    nil,
				handler: nil,
			},
			wantErr: true,
		},
		{
			name: "jwtProtectedInterceptor returns an error if no token provided in metadata",
			args: args{
				ctx:     ctxMdNoToken,
				req:     nil,
				info:    nil,
				handler: nil,
			},
			wantErr: true,
		},
		{
			name: "jwtProtectedInterceptor returns an error if invalid token provided in metadata",
			args: args{
				ctx:     ctxMdInvalidToken,
				req:     nil,
				info:    nil,
				handler: nil,
			},
			wantErr: true,
		},
		{
			name: "jwtProtectedInterceptor returns an error if user repository returns an error",
			args: args{
				ctx:     ctxMdValidToken,
				req:     nil,
				info:    nil,
				handler: nil,
			},
			wantErr: true,
		},
		{
			name: "jwtProtectedInterceptor returns an error if user repository returns an error",
			args: args{
				ctx:     ctxMdValidToken,
				req:     nil,
				info:    nil,
				handler: nil,
			},
			wantErr: true,
			before: func() {
				users.MakeUserRepository(users.MockRepository)
			},
			after: func() {
				users.MakeUserRepository(nil)
			},
		},
		{
			name: "jwtProtectedInterceptor calls the handler if valid user is found",
			args: args{
				ctx:     ctxMdValidToken,
				req:     nil,
				info:    nil,
				handler: testHandler,
			},
			wantResp: testUser,
			wantErr:  false,
			before: func() {
				users.MakeUserRepository(users.MockRepository)
				users.MockRepository.RetGetById = testUser
			},
			after: func() {
				users.MakeUserRepository(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			gotResp, err := jwtProtectedInterceptor(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if tt.after != nil {
				tt.after()
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtProtectedInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("jwtProtectedInterceptor() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_cspInterceptor(t *testing.T) {
	wantResp := "abc"
	testHandler := func(context.Context, any) (any, error) {
		return wantResp, nil
	}
	// note: experimental grpc API may break test in the future
	sts := &mockServerTransportStream{md: nil}
	ctxServerStream := grpc.NewContextWithServerTransportStream(context.Background(), sts)
	type args struct {
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name: "return error if used within non-grpc context",
			args: args{
				ctx:     context.Background(),
				req:     nil,
				info:    nil,
				handler: testHandler,
			},
			wantErr: true,
		},
		{
			name: "does not return error if used within non-grpc context",
			args: args{
				ctx:     ctxServerStream,
				req:     nil,
				info:    nil,
				handler: testHandler,
			},
			wantResp: wantResp,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := cspInterceptor(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("cspInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("cspInterceptor() = %v, want %v", gotResp, tt.wantResp)
				return
			}
			header, err := getHeader(*sts)
			if (err != nil) != tt.wantErr {
				t.Errorf("check csp header error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && header != config.AppConfig.ContentSecurityPolicy {
				t.Errorf("cspInterceptor header = %v, want %v", header, config.AppConfig.ContentSecurityPolicy)
			}
		})
	}
}
