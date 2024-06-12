package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/persistence"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserRepository interface {
	MakeUserHash(userID uuid.UUID) string
	CreateUser(ctx context.Context, userName string) (*model.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type UserRepositoryImpl struct {
	client       *redis.Client
	hashPrefix   string
	keyName      string
	keyCreatedAt string
}

var Repository UserRepository

var loggerRepository = basiclogger.BasicLogger{Namespace: "internal.users.repository"}

func init() {
	MakeUserRepository(nil)
}

func MakeUserRepository(r UserRepository) {
	if r != nil {
		Repository = r
	} else {
		Repository = &UserRepositoryImpl{
			persistence.GetClient(),
			"user",
			"name",
			"created_at",
		}
	}
}

// user:userUUID
func (r *UserRepositoryImpl) MakeUserHash(userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", r.hashPrefix, userID.String())
}

// Creates a new user
// Users are stored in Redis using the following structure
// [hash] => [key: value]
// [user:userUUID] => [name: userName] - Single record
// [user:userUUID] => [created_at: timestamppb.Timestamp.String()] - Single record
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, userName string) (*model.User, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	userID := uuid.New()

	pipe := r.client.Pipeline()
	userHash := r.MakeUserHash(userID)
	createdAt := timestamppb.Now()

	if _, err := pipe.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, userHash, r.keyName, userName)
		rdb.HSet(ctx, userHash, r.keyCreatedAt, createdAt.String())
		return nil
	}); err != nil {
		loggerRepository.LogError(reqID, "error creating user", "error", err)
		return nil, err
	}

	id := helpers.ProtoFromUUID(userID)
	u := model.User{
		Name:      userName,
		Id:        id,
		CreatedAt: createdAt,
	}

	loggerRepository.LogInfo(reqID, "user created and stored successfully", "user", &u)

	return &u, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	userHash := r.MakeUserHash(userID)
	userRaw, err := r.client.HGetAll(ctx, userHash).Result()

	if err != nil {
		loggerRepository.LogError(reqID, "error retrieving user", "error", err)
		return nil, err
	}

	name, ok := userRaw[r.keyName]

	if !ok || name == "" {
		loggerRepository.LogError(reqID, "required user name not provided")
		return nil, errors.New("required user name not provided")
	}

	createdAtRaw, ok := userRaw[r.keyCreatedAt]

	if !ok {
		loggerRepository.LogError(reqID, "error loading user: createdAt value not found")
		return nil, errors.New("error loading user")
	}

	var createdAt timestamppb.Timestamp
	err = prototext.Unmarshal([]byte(createdAtRaw), &createdAt)
	if err != nil {
		loggerRepository.LogError(reqID, "error un-marshaling create_at", "error", err)
		return nil, err
	}

	uid := helpers.ProtoFromUUID(userID)
	u := &model.User{
		Id:        uid,
		Name:      name,
		CreatedAt: &createdAt,
	}

	loggerRepository.LogInfo(reqID, "user fetched successfully", "user", u)

	return u, nil
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	reqID := helpers.GetReqIdFromContext(ctx)

	pipe := r.client.Pipeline()
	userHash := r.MakeUserHash(userID)

	_, err := r.GetByID(ctx, userID)
	if err != nil {
		loggerRepository.LogError(reqID, "error loading user", "error", err)
		return err
	}

	if _, err := pipe.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HDel(ctx, userHash, r.keyName)
		loggerRepository.LogDebug(reqID, "deleted user record", "userHash", userHash, "keyName", r.keyName)
		rdb.HDel(ctx, userHash, r.keyCreatedAt)
		loggerRepository.LogDebug(reqID, "deleted user record", "userHash", userHash, "keyCreatedAt", r.keyCreatedAt)
		return nil
	}); err != nil {
		loggerRepository.LogError(reqID, "error deleting user", "error", err)
		return err
	}

	loggerRepository.LogInfo(reqID, "user deleted successfully", "userID", userID)

	return nil
}
