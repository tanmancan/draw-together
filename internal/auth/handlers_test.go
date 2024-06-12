package auth

import (
	"context"
	"log/slog"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
)

func init() {
	level := slog.LevelError
	basiclogger.InitLogger(&level, nil)
}

func TestJwtAuthHandler(t *testing.T) {
	testUserID := uuid.New()
	testName := "test_name"
	ctx := context.Background()
	testToken, _ := GenerateUserToken(ctx, testUserID.String(), testName)
	testTokenInvalidUuid, _ := GenerateUserToken(ctx, "invalid_uuid", testName)
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantUid uuid.UUID
		wantErr bool
	}{
		{
			name: "JwtAuthHandler successfully parses token and returns user uuid",
			args: args{
				ctx:   ctx,
				token: testToken,
			},
			wantUid: testUserID,
			wantErr: false,
		},
		{
			name: "JwtAuthHandler returns error for empty token",
			args: args{
				ctx:   ctx,
				token: "",
			},
			wantErr: true,
		},
		{
			name: "JwtAuthHandler returns error for invalid token",
			args: args{
				ctx:   ctx,
				token: "invalid_token",
			},
			wantErr: true,
		},
		{
			name: "JwtAuthHandler returns error for token with invalid user ID",
			args: args{
				ctx:   ctx,
				token: testTokenInvalidUuid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUid, err := JwtAuthHandler(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("JwtAuthHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUid, tt.wantUid) {
				t.Errorf("JwtAuthHandler() = %v, want %v", gotUid, tt.wantUid)
			}
		})
	}
}
