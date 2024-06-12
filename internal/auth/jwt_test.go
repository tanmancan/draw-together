package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/config"
)

func TestParseUserToken(t *testing.T) {
	ctx := context.Background()
	testUserID := uuid.New()
	testName := uuid.NewString()
	token, _ := GenerateUserToken(ctx, testUserID.String(), testName)
	tinvaliduid := uuid.NewString()
	tinvalid := jwt.NewWithClaims(jwt.SigningMethodPS256,
		jwt.MapClaims{
			"id":   tinvaliduid,
			"name": "test_name",
		})
	key := config.AppConfig.JwtPublicKey
	tokenInvalidMethod, _ := tinvalid.SignedString([]byte(key))

	type args struct {
		ctx      context.Context
		tokenVal string
	}
	tests := []struct {
		name    string
		args    args
		want    jwt.MapClaims
		wantErr bool
	}{
		{
			name: "ParseUserToken successfully parses a valid token and returns claims",
			args: args{
				ctx:      ctx,
				tokenVal: token,
			},
			want: jwt.MapClaims{
				"id":   testUserID.String(),
				"name": testName,
			},
		},
		{
			name: "ParseUserToken returns error if token has invalid signing method",
			args: args{
				ctx:      ctx,
				tokenVal: tokenInvalidMethod,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUserToken(tt.args.ctx, tt.args.tokenVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUserToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateUserToken(t *testing.T) {
	ctx := context.Background()
	testUserId := uuid.New()
	testName := "test_name_2"
	type args struct {
		ctx      context.Context
		userID   string
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GenerateUserToken successfully generates and returns a token string",
			args: args{
				ctx:      ctx,
				userID:   testUserId.String(),
				userName: testName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateUserToken(tt.args.ctx, tt.args.userID, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GenerateUserToken() = %v, want %v", got, "string len > 0")
			}
		})
	}
}
