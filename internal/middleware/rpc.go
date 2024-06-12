package middleware

import (
	"context"

	"github.com/tanmancan/draw-together/internal/auth"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/service"
	"github.com/tanmancan/draw-together/internal/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type protectedMethodsInterceptors map[string]grpc.UnaryServerInterceptor

var loggerRpc = basiclogger.BasicLogger{Namespace: "internal.middleware.rpc"}

var protectedMethods = protectedMethodsInterceptors{
	// UserService
	service.UserService_CreateUser_FullMethodName: unprotectedInterceptor,
	service.UserService_GetUser_FullMethodName:    jwtProtectedInterceptor,
	service.UserService_DeleteUser_FullMethodName: jwtProtectedInterceptor,
	// PointerService
	service.PointerService_UpdatePointer_FullMethodName: jwtProtectedInterceptor,
	// ChatService
	service.ChatService_SendMessage_FullMethodName:      jwtProtectedInterceptor,
	service.ChatService_GetBoardMessages_FullMethodName: jwtProtectedInterceptor,
	// BoardService
	service.BoardService_CreateBoard_FullMethodName:      jwtProtectedInterceptor,
	service.BoardService_GetBoard_FullMethodName:         jwtProtectedInterceptor,
	service.BoardService_DeleteBoard_FullMethodName:      jwtProtectedInterceptor,
	service.BoardService_UpdateDrawing_FullMethodName:    jwtProtectedInterceptor,
	service.BoardService_GetBoardDrawings_FullMethodName: jwtProtectedInterceptor,
	service.BoardService_DrawingDetect_FullMethodName:    jwtProtectedInterceptor,
}

// Wrapper interceptor that calls specific handlers based on the current methods defined in protectedMethodsInterceptors
func authContextInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	wrapperInt, ok := protectedMethods[info.FullMethod]
	if !ok {
		err = status.Error(codes.NotFound, "not found")
		loggerRpc.LogError(reqID, "route or handler not found", "fullMethod", info.FullMethod, "error", err)
		return nil, err
	}

	loggerRpc.LogInfo(reqID, "found handler for method", "fullMethod", info.FullMethod)
	return wrapperInt(ctx, req, info, handler)
}

// Parses JWT token from metadata for a valid user ID
// If user is found, adds it to context
func jwtProtectedInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := status.Error(codes.InvalidArgument, "invalid metadata")
		loggerRpc.LogError(reqID, "error retrieving metadata", "error", err)
		return nil, err
	}

	var reqToken string = ""
	mdTokenVal := helpers.MetadataGetCi(md, helpers.HeaderXRequestToken)
	if len(mdTokenVal) > 0 {
		reqToken = mdTokenVal[0]
	}

	if len(reqToken) == 0 {
		err := status.Error(codes.PermissionDenied, "permission denied")
		loggerRpc.LogError(reqID, "request token not provided", "error", err)
		return nil, err
	}

	uid, err := auth.JwtAuthHandler(ctx, reqToken)
	if err != nil {
		err := status.Error(codes.PermissionDenied, "permission denied")
		loggerRpc.LogError(reqID, "error parsing request token", "error", err)
		return nil, err
	}

	user, err := users.Repository.GetByID(ctx, uid)
	if err != nil {
		err := status.Error(codes.PermissionDenied, "permission denied")
		loggerRpc.LogError(reqID, "error loading user", "error", err)
		return nil, err
	}

	if user == nil {
		err = status.Error(codes.PermissionDenied, "permission denied")
		loggerRpc.LogError(reqID, "user not found", "error", err)
		return nil, err
	}

	ctx = helpers.AddContextUser(ctx, user)
	loggerRpc.LogInfo(reqID, "user loaded successfully", "user", user)

	return handler(ctx, req)
}

// Unprotected method interceptor - will call handler without any authentication
func unprotectedInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	loggerRpc.LogInfo(reqID, "handling unprotected route", "fullMethod", info.FullMethod)
	return handler(ctx, req)
}

// Adds tracking requestID to context and grpc response header
func requestIdInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx = helpers.CreateContextWithRequestID(ctx)
	reqID := helpers.GetReqIdFromContext(ctx)
	md := metadata.MD{}
	md.Set(helpers.HeaderXRequestUuid, reqID)
	grpc.SetHeader(ctx, md)
	return handler(ctx, req)
}

// Adds Content Security Policy header
func cspInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	headers := metadata.Pairs(helpers.HeaderContentSecurityPolicy, config.AppConfig.ContentSecurityPolicy)
	err = grpc.SetHeader(ctx, headers)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

// Logs every request
func loggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	loggerRpc.LogInfo(reqID, "grpc request", "fullMethod", info.FullMethod)
	return handler(ctx, req)
}

func GetChainUnaryInterceptor() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		requestIdInterceptor,
		authContextInterceptor,
		cspInterceptor,
		loggerInterceptor,
	)
}
