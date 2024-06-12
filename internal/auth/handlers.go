package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"google.golang.org/grpc/codes"
)

var loggerHandlers = basiclogger.BasicLogger{Namespace: "internal.auth.handlers"}

// Validates user JWT token and returns user ID if valid
func JwtAuthHandler(ctx context.Context, token string) (uid uuid.UUID, err error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	if token == "" {
		err := helpers.MakeError(ctx, codes.PermissionDenied, "request token not provided")
		loggerHandlers.LogError(reqID, "error retrieving auth token", "error", err)

		return uuid.Nil, err
	}

	loggerHandlers.LogInfo(reqID, "parsing token start")
	claims, err := parseUserToken(ctx, token)
	if err != nil {
		return uuid.Nil, err
	}

	id := claims["id"].(string)
	uid, err = uuid.Parse(id)
	if err != nil {
		errParse := helpers.MakeError(ctx, codes.Unauthenticated, "invalid token")
		loggerHandlers.LogError(reqID, "error parsing uuid in claims", "error", err, "id", claims["id"])
		return uuid.Nil, errParse
	}

	loggerHandlers.LogInfo(reqID, "parsing token successful")

	return uid, nil
}
