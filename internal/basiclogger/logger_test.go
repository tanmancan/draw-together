package basiclogger

import (
	"bytes"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestInitLogger(t *testing.T) {
	testUUID := uuid.NewString()
	type args struct {
		level slog.Level
	}
	tests := []struct {
		name     string
		args     args
		debugLog bool
		testFunc func()
		wantW    string
	}{
		{
			name: "InitLogger sets the correct log level",
			args: args{
				level: slog.LevelInfo,
			},
			testFunc: func() {
				testLogger := BasicLogger{Namespace: "test"}
				testLogger.LogInfo("test", testUUID)
			},
			wantW: testUUID,
		},
		{
			name: "InitLogger sets the correct log level",
			args: args{
				level: slog.LevelWarn,
			},
			testFunc: func() {
				testLogger := BasicLogger{Namespace: "test"}
				testLogger.LogWarn("test", testUUID)
			},
			wantW: testUUID,
		},
		{
			name: "InitLogger sets the correct log level",
			args: args{
				level: slog.LevelError,
			},
			testFunc: func() {
				testLogger := BasicLogger{Namespace: "test"}
				testLogger.LogError("test", testUUID)
			},
			wantW: testUUID,
		},
		{
			name: "InitLogger won't log message below specified level",
			args: args{
				level: slog.LevelWarn,
			},
			testFunc: func() {
				testLogger := BasicLogger{Namespace: "test"}
				testLogger.LogInfo("test", testUUID)
			},
			wantW: "",
		},
		{
			name: "InitLogger does not log debug messages unless 'debuglog' flag is set",
			args: args{
				level: slog.LevelDebug,
			},
			testFunc: func() {
				testLogger := BasicLogger{Namespace: "test"}
				testLogger.LogDebug("test", testUUID)
			},
			wantW: "",
		},
		{
			name: "InitLogger log debug messages when debug log is set",
			args: args{
				level: slog.LevelDebug,
			},
			debugLog: true,
			testFunc: func() {
				testLogger := BasicLogger{Namespace: "test"}
				testLogger.LogDebug("test", testUUID)
			},
			wantW: testUUID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			debuglog = tt.debugLog
			InitLogger(&tt.args.level, w)
			tt.testFunc()
			if gotW := w.String(); !strings.Contains(gotW, tt.wantW) {
				t.Errorf("InitLogger() = %v, want %v", gotW, tt.wantW)
			}
			debuglog = false
		})
	}
}

func Test_replaceAttr(t *testing.T) {
	method := "POST"
	host := "localhost/test1"
	origin := "localhost/origin"

	req, _ := http.NewRequest(method, host, nil)
	req.Header.Set("Origin", origin)

	user := &model.User{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		Name: uuid.NewString(),
	}

	board := &model.Board{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		Name:      uuid.NewString(),
		CreatedAt: timestamppb.Now(),
		Owner:     user,
	}

	type args struct {
		groups []string
		a      slog.Attr
	}
	tests := []struct {
		name string
		args args
		want slog.Attr
	}{
		{
			name: "replaceAttr prints method, host, and origin fields for http.Request",
			args: args{
				a: slog.Any("req", req),
			},
			want: slog.Group("request", "method", method, "host", host, "origin", origin),
		},
		{
			name: "replaceAttr prints formatted values for model.Board",
			args: args{
				a: slog.Any("b", board),
			},
			want: slog.Group("board",
				"id", board.Id.Value,
				"name", board.Name,
				"createdAt", board.CreatedAt.AsTime(),
				"owner", slog.Group("owner",
					"id", board.Owner.Id.Value,
					"name", board.Owner.Name,
					"createdAt", board.Owner.CreatedAt.AsTime())),
		},
		{
			name: "replaceAttr prints formatted values for model.User",
			args: args{
				a: slog.Any("u", user),
			},
			want: slog.Group("user",
				"id", user.Id.Value,
				"name", user.Name,
				"createdAt", user.CreatedAt.AsTime()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceAttr(tt.args.groups, tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replaceAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
