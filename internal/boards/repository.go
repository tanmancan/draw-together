package boards

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/persistence"
	"github.com/tanmancan/draw-together/internal/users"
	"google.golang.org/protobuf/encoding/prototext"
)

// Map<userUuid, boardUuid>
type BoardUsers map[uuid.UUID]uuid.UUID

// Map<userUuid, boardUuid>
type BoardOwners map[uuid.UUID]uuid.UUID

// Map<boardUUID, isBoardOwner>
type UserBoards map[uuid.UUID]bool

// Map<string(userUUID), ImageData>
type BoardDrawings map[string]*model.ImageData

type BoardRepository interface {
	CreateBoard(ctx context.Context, name string, owner *model.User) (*model.Board, error)
	DeleteBoard(ctx context.Context, boardID uuid.UUID) error
	GetByID(ctx context.Context, boardID uuid.UUID) (*model.Board, error)
	GetBoardsByUser(ctx context.Context, userID uuid.UUID) ([]*model.Board, error)
	GetAllUsers(ctx context.Context, boardID uuid.UUID) ([]*model.User, error)
	AddUser(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	BoardHasUser(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) bool
	GetUserBoards(ctx context.Context, userID uuid.UUID) (UserBoards, error)
	GetBoardDrawings(ctx context.Context, boardID uuid.UUID) (BoardDrawings, error)
	UpdateBoardDrawing(ctx context.Context, boardID uuid.UUID, userID uuid.UUID, imageData *model.ImageData) error
}

type BoardRepositoryImpl struct {
	client        *redis.Client
	boardPrefix   string
	userPrefix    string
	drawingPrefix string
	keyProto      string
	maxUsers      uint8
}

var Repository BoardRepository

var loggerRepository = basiclogger.BasicLogger{Namespace: "internal.boards.repository"}

func init() {
	MakeBoardRepository(nil)
}

func MakeBoardRepository(r BoardRepository) {
	if r != nil {
		Repository = r
	} else {
		Repository = &BoardRepositoryImpl{
			persistence.GetClient(),
			"board",
			"user",
			"drawing",
			"proto",
			5,
		}
	}
}

// Board object hash: board:boardUUID
// Stores data for an individual board
func (r *BoardRepositoryImpl) MakeBoardObjectHash(boardID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", r.boardPrefix, boardID.String())
}

// Board users hash: board:user:boardUUID
// Stores list of users registered to a board
func (r *BoardRepositoryImpl) MakeBoardUsersHash(boardID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", r.boardPrefix, r.userPrefix, boardID.String())
}

// User boards hash: user:board:userUUID
// Stores list of boards joined by a user
func (r *BoardRepositoryImpl) MakeUserBoardsHash(userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", r.userPrefix, r.boardPrefix, userID.String())
}

// Board drawing hash: board:drawing:boardUUID
// Stores a list of user drawing layers for a board
func (r *BoardRepositoryImpl) MakeBoardDrawingHash(boardID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", r.boardPrefix, r.drawingPrefix, boardID.String())
}

// Creates a new board
// Boards are stored in Redis using the following structure
// [hash] => [key: value]
// [board:boardUUID] => [proto: model.Board.String()] - Single record
// [board:users:boardUUID] => [userUUID: boardUUID] - Multiple record
func (r *BoardRepositoryImpl) CreateBoard(ctx context.Context, name string, owner *model.User) (*model.Board, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	newBoardUuid := uuid.New()
	ownerID, err := helpers.ProtoToUUID(owner.Id)

	if err != nil {
		loggerRepository.LogError(reqID, "error parsing owner_id", "error", err)
		return nil, err
	}

	var board *model.Board = &model.Board{
		Id: &model.UUID{
			Value: newBoardUuid.String(),
		},
		Name:       name,
		Owner:      owner,
		BoardUsers: nil,
	}

	pipe := r.client.Pipeline()
	boardHash := r.MakeBoardObjectHash(newBoardUuid)
	usersHash := r.MakeBoardUsersHash(newBoardUuid)
	userBoardHash := r.MakeUserBoardsHash(ownerID)

	if _, err := pipe.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, boardHash, r.keyProto, board.String())
		rdb.HSet(ctx, usersHash, ownerID.String(), newBoardUuid.String())
		rdb.HSet(ctx, userBoardHash, ownerID.String(), true)
		return nil
	}); err != nil {
		loggerRepository.LogError(reqID, "error creating board", "error", err)
		return nil, err
	}

	loggerRepository.LogDebug(reqID, "board created and stored", "board", board)

	return board, nil
}

func (r *BoardRepositoryImpl) GetByID(ctx context.Context, boardID uuid.UUID) (*model.Board, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	boardHash := r.MakeBoardObjectHash(boardID)
	boardProtoStr, err := r.client.HGet(ctx, boardHash, r.keyProto).Result()

	if err != nil {
		loggerRepository.LogError(reqID, "error retrieving board", "error", err)
		return nil, err
	}

	var board *model.Board = &model.Board{}
	err = prototext.Unmarshal([]byte(boardProtoStr), board)
	if err != nil {
		loggerRepository.LogError(reqID, "error parsing board data", "error", err)
		return nil, err
	}

	user := helpers.GetUserFromContext(ctx)
	if user == nil {
		loggerRepository.LogError(reqID, "user not found")
		return nil, errors.New("user not found")
	}

	boardUsers, err := r.GetAllUsers(ctx, boardID)

	if err != nil {
		loggerRepository.LogError(reqID, "error loading board users: ", "boardID", boardID, "error", err)
		return nil, errors.New("board not found")
	}

	board.BoardUsers = boardUsers
	loggerRepository.LogDebug(reqID, "found board", "board", board)

	return board, nil
}

// Adds a new user to a board
// Stored in Redis using the following structure
// [hash] => [key: value]
// [board:users:boardUUID] => [userUUID: boardUUID] - Multiple record
func (r *BoardRepositoryImpl) GetAllUsers(ctx context.Context, boardID uuid.UUID) ([]*model.User, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	usersHash := r.MakeBoardUsersHash(boardID)
	boardUsers, err := r.client.HGetAll(ctx, usersHash).Result()

	if err != nil {
		loggerRepository.LogError(reqID, "error retrieving all users", "error", err)
		return nil, err
	}

	var ulist []*model.User

	for uidStr := range boardUsers {
		uid, err := uuid.Parse(uidStr)

		if err != nil {
			loggerRepository.LogError(reqID, "invalid user uuid", "userID", uid, "error", err)
			continue
		}

		u, err := users.Repository.GetByID(ctx, uid)

		if err != nil {
			loggerRepository.LogError(reqID, "user not found for requested board", "boardID", boardID, "userID", uid, "error", err)
			continue
		}

		ulist = append(ulist, u)
	}

	return ulist, nil

}

func (r *BoardRepositoryImpl) AddUser(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error {
	reqID := helpers.GetReqIdFromContext(ctx)

	if ok := r.BoardHasUser(ctx, boardID, userID); ok {
		loggerRepository.LogWarn(reqID, "user already exist in board: noop")
		return nil
	}

	usersHash := r.MakeBoardUsersHash(boardID)
	userBoardHash := r.MakeUserBoardsHash(userID)
	boardUsers, _ := r.GetAllUsers(ctx, boardID)
	uc := len(boardUsers)

	if uc >= int(r.maxUsers) {
		err := errors.New("max user limit reached for board")
		loggerRepository.LogError(reqID, "max user limit reached for board", "userCount", uc, "boardID", boardID, "error", err)
		return err
	}

	pipe := r.client.Pipeline()

	if _, err := pipe.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, usersHash, userID.String(), boardID.String())
		rdb.HSet(ctx, userBoardHash, boardID.String(), false)
		return nil
	}); err != nil {
		loggerRepository.LogError(reqID, "error adding user to board", "error", err)
		return err
	}

	return nil
}

func (r *BoardRepositoryImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	userBoards, err := r.GetUserBoards(ctx, userID)
	loggerRepository.LogInfo(reqID, "delete board user started", "userID", userID)
	if err != nil {
		loggerRepository.LogError(reqID, "error loading user boards", "error", err)
		return err
	}

	for boardID, isOwner := range userBoards {
		if isOwner {
			loggerRepository.LogInfo(reqID, "user-owned board found")
			err := r.DeleteBoard(ctx, boardID)
			if err != nil {
				loggerRepository.LogError(reqID, "error deleting user owned board", "error", err)
				continue
			}
			loggerRepository.LogInfo(reqID, "user-owned board deleted")
		}
	}

	userBoardHash := r.MakeUserBoardsHash(userID)
	r.client.HDel(ctx, userBoardHash)

	loggerRepository.LogInfo(reqID, "delete board user completed", "userID", userID)

	return nil
}

func (r *BoardRepositoryImpl) DeleteBoard(ctx context.Context, boardID uuid.UUID) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	board, err := r.GetByID(ctx, boardID)

	if err != nil {
		loggerRepository.LogError(reqID, "error retrieving board", "error", err)
		return err
	}

	if board != nil {
		err = errors.New("board not found")
		loggerRepository.LogError(reqID, "error retrieving board", "error", err)
		return err
	}

	pipe := r.client.Pipeline()
	boardHash := r.MakeBoardObjectHash(boardID)
	usersHash := r.MakeBoardUsersHash(boardID)
	boardUsers, err := r.GetAllUsers(ctx, boardID)

	if err != nil {
		loggerRepository.LogError(reqID, "error retrieving all users", "error", err)
		return err
	}

	if _, err := pipe.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HDel(ctx, boardHash, boardID.String())
		loggerRepository.LogDebug(reqID, "deleted board", "boardHash", boardHash, "boardID", boardID.String())

		for _, user := range boardUsers {
			userID, _ := helpers.ProtoToUUID(user.GetId())
			userBoardHash := r.MakeUserBoardsHash(userID)
			rdb.HDel(ctx, usersHash, user.Id.Value)
			loggerRepository.LogInfo(reqID, "deleted board", "userHash", usersHash, "userID", user.Id)
			rdb.HDel(ctx, userBoardHash, board.Id.Value)
			loggerRepository.LogInfo(reqID, "deleted board user", "userBoardHash", userBoardHash, "userID", user.Id)
		}

		return nil
	}); err != nil {
		loggerRepository.LogError(reqID, "error deleting board", "error", err)
		return err
	}

	return nil
}

func (r *BoardRepositoryImpl) BoardHasUser(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) bool {
	reqID := helpers.GetReqIdFromContext(ctx)
	usersHash := r.MakeBoardUsersHash(boardID)
	bid, err := r.client.HGet(ctx, usersHash, userID.String()).Result()

	if err != nil {
		loggerRepository.LogWarn(reqID, "board user not found", "userID", userID, "boardID", boardID, "error", err)
	}

	return bid == boardID.String()
}

func (r *BoardRepositoryImpl) GetBoardsByUser(ctx context.Context, userID uuid.UUID) ([]*model.Board, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	userBoards, err := r.GetUserBoards(ctx, userID)

	if err != nil {
		loggerRepository.LogError(reqID, "error loading user boards", "error", err)
		return nil, err
	}

	var boards []*model.Board

	for boardID := range userBoards {
		board, err := r.GetByID(ctx, boardID)

		if err != nil {
			loggerRepository.LogError(reqID, "error loading board", "error", err)
			return nil, err
		}

		boards = append(boards, board)
	}

	return boards, nil
}

func (r *BoardRepositoryImpl) GetUserBoards(ctx context.Context, userID uuid.UUID) (UserBoards, error) {
	var userBoards UserBoards = make(UserBoards)
	reqID := helpers.GetReqIdFromContext(ctx)
	userBoardHash := r.MakeUserBoardsHash(userID)
	boardIds, err := r.client.HGetAll(ctx, userBoardHash).Result()

	if err != nil {
		loggerRepository.LogError(reqID, "error loading user boards", "error", err)
	}

	for bid, o := range boardIds {
		boardID, err := uuid.Parse(bid)

		if err != nil {
			loggerRepository.LogError(reqID, "error parsing user board ID", "error", err)
			return nil, err
		}

		isOwner, err := strconv.ParseBool(o)

		if err != nil {
			loggerRepository.LogError(reqID, "error parsing user board owner value", "error", err)
			return nil, err
		}

		_, err = r.GetByID(ctx, boardID)

		if err != nil {
			loggerRepository.LogError(reqID, "error loading board", "error", err)
			return nil, err
		}

		userBoards[boardID] = isOwner
	}

	return userBoards, nil
}

func (r *BoardRepositoryImpl) UpdateBoardDrawing(ctx context.Context, boardID uuid.UUID, userID uuid.UUID, imageData *model.ImageData) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	if ok := r.BoardHasUser(ctx, boardID, userID); !ok {
		loggerRepository.LogError(reqID, "user does not belong to board", "userID", userID, "boardID", boardID)
		return errors.New("board not found")
	}

	drawingHash := r.MakeBoardDrawingHash(boardID)
	rdb := persistence.GetClient()

	loggerRepository.LogInfo(reqID, "saving drawing: start", "userID", userID, "boardID", boardID)
	err := rdb.HSet(ctx, drawingHash, userID.String(), imageData.String()).Err()
	if err != nil {
		loggerRepository.LogError(reqID, "error saving drawing", "userID", userID, "boardID", boardID)
	}

	loggerRepository.LogInfo(reqID, "saving drawing: done", "userID", userID, "boardID", boardID)

	return nil
}

func (r *BoardRepositoryImpl) GetBoardDrawings(ctx context.Context, boardID uuid.UUID) (BoardDrawings, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	rdb := persistence.GetClient()
	drawingHash := r.MakeBoardDrawingHash(boardID)

	loggerRepository.LogInfo(reqID, "get image data for board start", "boardID", boardID)
	drawings, err := rdb.HGetAll(ctx, drawingHash).Result()
	if err != nil {
		loggerRepository.LogError(reqID, "error loading board drawings", "boardID", boardID, "error", err)
		return nil, err
	}
	var boardDrawings BoardDrawings = make(BoardDrawings)

	for userID, imgData := range drawings {
		var iData model.ImageData
		err = prototext.Unmarshal([]byte(imgData), &iData)
		if err != nil {
			loggerRepository.LogError(reqID, "error parsing image data", "userID", userID, "boardID", boardID)
			continue
		}

		loggerRepository.LogInfo(reqID, "found image data for board user", "userID", userID, "boardID", boardID)
		boardDrawings[userID] = &iData
	}

	return boardDrawings, nil
}
