package basiclogger

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/model"
)

type BasicLogger struct {
	Namespace string
}

type key int

const (
	ReqID key = iota
)

var debuglog bool = config.AppConfig.LoggerConfig.EnableDebug

var logLevel slog.Level = slog.LevelInfo

var slogger *slog.Logger

var logger = BasicLogger{Namespace: "internal.basiclogger"}

func InitLogger(level *slog.Level, w io.Writer) {
	if level != nil {
		logLevel = *level
	}

	if w == nil {
		w = os.Stdout
	}

	slogger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel,
		ReplaceAttr: replaceAttr,
	}))

	if debuglog {
		logger.LogDebug(
			uuid.NewString(),
			"Debug log is enabled",
		)
	}
}

func init() {
	if debuglog {
		logLevel = slog.LevelDebug
	}
	InitLogger(&logLevel, nil)
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
		source.Function = filepath.Base(source.Function)
	}

	switch a.Value.Any().(type) {
	case *http.Request:
		req := a.Value.Any().(*http.Request)
		a = slog.Group("request", "method", req.Method, "host", req.Host, "origin", req.Header.Get("origin"))
	case *model.Board:
		b := a.Value.Any().(*model.Board)
		a = slog.Group("board",
			"id", b.Id.Value,
			"name", b.Name,
			"createdAt", b.CreatedAt.AsTime(),
			"owner", slog.Group("owner",
				"id", b.Owner.Id.Value,
				"name", b.Owner.Name,
				"createdAt", b.Owner.CreatedAt.AsTime()))
	case *model.User:
		u := a.Value.Any().(*model.User)
		a = slog.Group("user",
			"id", u.Id.Value,
			"name", u.Name,
			"createdAt", u.CreatedAt.AsTime())
	default:
	}

	return a
}

func (l *BasicLogger) logInternal(level slog.Level, reqID string, msg string, supplemental ...interface{}) {
	ctx := context.WithValue(context.Background(), ReqID, reqID)
	if !slogger.Handler().Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // Skip [Callers, logInternal, Log[Level]]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])

	r.Add("reqID", reqID)
	r.Add("namespace", l.Namespace)
	r.Add(supplemental...)
	_ = slogger.Handler().Handle(ctx, r)
}

func (l *BasicLogger) LogDebug(reqID string, msg string, supplemental ...interface{}) {
	if debuglog {
		l.logInternal(slog.LevelDebug, reqID, msg, supplemental...)
	}
}

func (l *BasicLogger) LogInfo(reqID string, msg string, supplemental ...interface{}) {
	l.logInternal(slog.LevelInfo, reqID, msg, supplemental...)
}

func (l *BasicLogger) LogWarn(reqID string, msg string, supplemental ...interface{}) {
	l.logInternal(slog.LevelWarn, reqID, msg, supplemental...)
}

func (l *BasicLogger) LogError(reqID string, msg string, supplemental ...interface{}) {
	l.logInternal(slog.LevelError, reqID, msg, supplemental...)
}

func (l *BasicLogger) LogFatal(reqID string, msg string, supplemental ...interface{}) {
	l.logInternal(slog.LevelError, reqID, msg, supplemental...)
	os.Exit(1)
}
