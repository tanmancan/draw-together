package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DefaultBoardConfigType struct {
	MaxUser     int
	MaxLifetime time.Duration
	MaxGuesses  int
}

type NetworkConfigType struct {
	Host string
	Port int
}

type WebsocketConfigType struct {
	NetworkConfigType
	ReadDeadline  time.Duration
	WriteDeadline time.Duration
	PingInterval  time.Duration
	ReadBuffer    int
	WriteBuffer   int
	PongDeadline  int
	ReadLimit     int64
}

type ChannelConfigType struct {
	Name string
	Id   string
}

type EventConfigType struct {
	ChannelChatMessage   ChannelConfigType
	ChannelUpdatePointer ChannelConfigType
	ChannelUpdateDrawing ChannelConfigType
	ChannelUpdateBoard   ChannelConfigType
	ChannelUpdateUser    ChannelConfigType
	ChannelDetectDrawing ChannelConfigType
	ChannelDetectQueue   ChannelConfigType
}

type DetectionConfigType struct {
	AzureCvEndpoint string
	AzureCvKey      string
}

type LoggerConfigType struct {
	EnableDebug bool
}

type AllowedOriginList []string

type CorsConfigType struct {
	AllowedOrigin AllowedOriginList
}

type AppConfigType struct {
	DefaultBoardConfig        DefaultBoardConfigType
	NetworkServiceGrpc        NetworkConfigType
	NetworkServiceWebsocket   WebsocketConfigType
	NetworkServiceRedisClient NetworkConfigType
	EventConfig               EventConfigType
	DetectionConfig           DetectionConfigType
	CorsConfig                CorsConfigType
	LoggerConfig              LoggerConfigType
	JwtPublicKey              string
	ContentSecurityPolicy     string
}

const (
	SERVICE_GRPC_HOST      string = "SERVICE_GRPC_HOST"
	SERVICE_GRPC_PORT      string = "SERVICE_GRPC_PORT"
	SERVICE_WEBSOCKET_HOST string = "SERVICE_WEBSOCKET_HOST"
	SERVICE_WEBSOCKET_PORT string = "SERVICE_WEBSOCKET_PORT"
	SERVICE_REDIS          string = "SERVICE_REDIS"
	REDIS_PORT             string = "REDIS_PORT"
	CORS_ALLOWED_ORIGIN    string = "CORS_ALLOWED_ORIGIN"
	JWT_PUBLIC_KEY         string = "JWT_PUBLIC_KEY"
	LOGGER_DEBUG           string = "LOGGER_DEBUG"
	AZURE_CV_ENDPOINT      string = "AZURE_CV_ENDPOINT"
	AZURE_CV_KEY           string = "AZURE_CV_KEY"
)

var showVersion bool

var BuildVersion string

var GitSHA string

var BuildTime string

var AppConfig AppConfigType

func defaultBoardConfig() DefaultBoardConfigType {
	return DefaultBoardConfigType{
		MaxUser:     5,
		MaxLifetime: time.Minute * 5,
		MaxGuesses:  5,
	}
}

func grpcConfig() NetworkConfigType {
	host := getEnvStr(SERVICE_GRPC_HOST, "0.0.0.0")
	port := getEnvInt(SERVICE_GRPC_PORT, 9090)
	return NetworkConfigType{
		Host: host,
		Port: port,
	}
}

func websocketConfig() WebsocketConfigType {
	host := getEnvStr(SERVICE_WEBSOCKET_HOST, "0.0.0.0")
	port := getEnvInt(SERVICE_WEBSOCKET_PORT, 9003)
	return WebsocketConfigType{
		NetworkConfigType: NetworkConfigType{
			Host: host,
			Port: port,
		},
		ReadLimit:     512,
		ReadDeadline:  30 * time.Second,
		WriteDeadline: 5 * time.Second,
		PingInterval:  25 * time.Second,
	}
}

func redisClientConfig() NetworkConfigType {
	host := getEnvStr(SERVICE_REDIS, "localhost")
	port := getEnvInt(REDIS_PORT, 6379)
	return NetworkConfigType{
		Host: host,
		Port: port,
	}
}

func eventConfig() EventConfigType {
	return EventConfigType{
		ChannelChatMessage: ChannelConfigType{
			Name: "chat",
			Id:   uuid.NewString(),
		},
		ChannelUpdatePointer: ChannelConfigType{
			Name: "pointer",
			Id:   uuid.NewString(),
		},
		ChannelUpdateDrawing: ChannelConfigType{
			Name: "drawing",
			Id:   uuid.NewString(),
		},
		ChannelUpdateBoard: ChannelConfigType{
			Name: "board",
			Id:   uuid.NewString(),
		},
		ChannelUpdateUser: ChannelConfigType{
			Name: "user",
			Id:   uuid.NewString(),
		},
		ChannelDetectDrawing: ChannelConfigType{
			Name: "detect",
			Id:   uuid.NewString(),
		},
		ChannelDetectQueue: ChannelConfigType{
			Name: "detect_queue",
			Id:   uuid.NewString(),
		},
	}
}

func detectConfig() DetectionConfigType {
	azureCvEndpoint := getEnvStr(AZURE_CV_ENDPOINT, "")
	azureCvKey := getEnvStr(AZURE_CV_KEY, "")
	return DetectionConfigType{
		AzureCvEndpoint: azureCvEndpoint,
		AzureCvKey:      azureCvKey,
	}
}

func corsConfig() CorsConfigType {
	allowedList := strings.Split(getEnvStr(CORS_ALLOWED_ORIGIN, "http://localhost:5173"), ",")
	return CorsConfigType{
		AllowedOrigin: append(AllowedOriginList{}, allowedList...),
	}
}

func loggerConfig() LoggerConfigType {
	enableDebug := getEnvBool(LOGGER_DEBUG, false)
	return LoggerConfigType{
		EnableDebug: enableDebug,
	}
}

func init() {
	AppConfig = AppConfigType{
		DefaultBoardConfig:        defaultBoardConfig(),
		NetworkServiceGrpc:        grpcConfig(),
		NetworkServiceWebsocket:   websocketConfig(),
		NetworkServiceRedisClient: redisClientConfig(),
		EventConfig:               eventConfig(),
		DetectionConfig:           detectConfig(),
		CorsConfig:                corsConfig(),
		JwtPublicKey:              getEnvStr(JWT_PUBLIC_KEY, ""),
		LoggerConfig:              loggerConfig(),
		ContentSecurityPolicy:     "default-src 'none';frame-ancestors 'none';",
	}
}

func SetupConfigFlags() {
	flag.BoolVar(&showVersion, "version", false, "print version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %s\nSHA: %s\nBuildTime: %s\n", BuildVersion, GitSHA, BuildTime)
		os.Exit(1)
	}
}
