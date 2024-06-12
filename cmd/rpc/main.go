package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/boards"
	"github.com/tanmancan/draw-together/internal/chat"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/creds"
	"github.com/tanmancan/draw-together/internal/middleware"
	"github.com/tanmancan/draw-together/internal/pointer"
	"github.com/tanmancan/draw-together/internal/service"
	"github.com/tanmancan/draw-together/internal/users"
	"google.golang.org/grpc"
)

var logger = basiclogger.BasicLogger{Namespace: "rpc.main"}

func main() {
	config.SetupConfigFlags()
	reqID := uuid.NewString()
	h := config.AppConfig.NetworkServiceGrpc.Host
	p := config.AppConfig.NetworkServiceGrpc.Port

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", h, p))
	if err != nil {
		logger.LogError(reqID, "error starting grpc service", "error", err)
	}

	msg := fmt.Sprintf("grpc service listening at %s:%d", h, p)
	logger.LogInfo(reqID, msg)

	var opts []grpc.ServerOption

	chainUnaryInt := middleware.GetChainUnaryInterceptor()
	c, err := creds.GetTransportCredentials()
	if err != nil {
		log.Fatal(err)
	}
	creds := grpc.Creds(c)
	opts = append(opts, creds, chainUnaryInt)
	s := grpc.NewServer(opts...)

	service.RegisterUserServiceServer(s, &users.UserServiceServerImpl{})
	service.RegisterBoardServiceServer(s, &boards.BoardServiceServerImpl{})
	service.RegisterChatServiceServer(s, &chat.ChatServiceServerImpl{})
	service.RegisterPointerServiceServer(s, &pointer.PointerServiceServerImpl{})

	i := s.GetServiceInfo()
	logger.LogInfo(reqID, "service info", "info", i)

	log.Fatal(s.Serve(lis), reqID)
}
