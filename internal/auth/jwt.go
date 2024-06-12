package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"google.golang.org/grpc/codes"
)

var loggerJwt = basiclogger.BasicLogger{Namespace: "internal.auth.jwt"}

// Generate JWT using user ID and user name as claims
func GenerateUserToken(ctx context.Context, userID string, userName string) (string, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	loggerJwt.LogInfo(reqID, "generate user token: start")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   userID,
			"name": userName,
		})

	key := config.AppConfig.JwtPublicKey
	str, err := t.SignedString([]byte(key))

	if err != nil {
		errSign := helpers.MakeError(ctx, codes.Internal, "error generating user token")
		loggerJwt.LogError(reqID, "generating user token", "error", err)
		return "", errSign
	}

	loggerJwt.LogInfo(reqID, "generate user token: done")

	return str, nil
}

// Parse JWT user token and returns claims if successful
func parseUserToken(ctx context.Context, tokenVal string) (jwt.MapClaims, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	token, err := jwt.Parse(tokenVal, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errInvalidMethod := helpers.MakeError(ctx, codes.Unauthenticated, "unexpected signing method")
			loggerJwt.LogError(reqID, "unexpected signing method", "alg", token.Header["alg"])
			return nil, errInvalidMethod
		}

		key := config.AppConfig.JwtPublicKey
		return []byte(key), nil
	})

	if err != nil {
		errParse := helpers.MakeError(ctx, codes.Unauthenticated, "invalid token")
		loggerJwt.LogError(reqID, "error parse user token", "error", err, "errPrase", errParse)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		loggerJwt.LogInfo(reqID, "parse user token successful", "claims", claims)
		return claims, nil
	} else {
		errParse := helpers.MakeError(ctx, codes.Unauthenticated, "invalid token")
		loggerJwt.LogError(reqID, "parse user token error: unable to retrieve claims.")
		return nil, errParse
	}
}
