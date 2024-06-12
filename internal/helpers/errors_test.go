package helpers

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestMakeError(t *testing.T) {
	ctx := CreateContextWithRequestID(context.Background())
	reqID := GetReqIdFromContext(ctx)
	type args struct {
		ctx  context.Context
		c    codes.Code
		msg  string
		args []any
	}
	tests := []struct {
		name       string
		args       args
		wantErrStr string
	}{
		{
			name: "generates errors that contains the provided message",
			args: args{
				ctx: ctx,
				c:   codes.Aborted,
				msg: "error test",
			},
			wantErrStr: "error test",
		},
		{
			name: "generates errors message that contain the reqID from passed context",
			args: args{
				ctx: ctx,
				c:   codes.Aborted,
				msg: "error test",
			},
			wantErrStr: reqID,
		},
		{
			name: "generates that contain the provided code",
			args: args{
				ctx: ctx,
				c:   codes.Aborted,
				msg: "error test",
			},
			wantErrStr: "code = Aborted",
		},
		{
			name: "generates that contain the provided code",
			args: args{
				ctx: ctx,
				c:   codes.NotFound,
				msg: "error test",
			},
			wantErrStr: "code = NotFound",
		},
		{
			name: "generates error that contain the error code prefix",
			args: args{
				ctx: ctx,
				c:   codes.NotFound,
				msg: "error test",
			},
			wantErrStr: "not found: ",
		},
		{
			name: "generates error that contain the error code prefix",
			args: args{
				ctx: ctx,
				c:   codes.AlreadyExists,
				msg: "error test",
			},
			wantErrStr: "already exists: ",
		},
		{
			name: "generates message using provided args",
			args: args{
				ctx:  ctx,
				c:    codes.NotFound,
				msg:  "error test %s",
				args: []any{"test_arg"},
			},
			wantErrStr: "test_arg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MakeError(tt.args.ctx, tt.args.c, tt.args.msg, tt.args.args...); err != nil {
				msg := err.Error()
				if ok := strings.Contains(msg, tt.wantErrStr); !ok {
					t.Errorf("MakeError() error = %v, wantErrStr contains %v", err, tt.wantErrStr)
				}
			}
		})
	}
}
