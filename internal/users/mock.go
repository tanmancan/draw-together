package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/model"
)

type MockUserRepository struct {
	CallMakeUserHashCt int
	CallCreateUserCt   int
	CallGetByIDCt      int
	CallDeleteUserCt   int
	RetMakeUserHash    string
	RetCreateUser      *model.User
	RetGetById         *model.User
	ErrCreateUser      error
	ErrGetById         error
	ErrDeleteUser      error
}

func (r *MockUserRepository) MakeUserHash(userID uuid.UUID) string {
	r.CallMakeUserHashCt++
	return r.RetMakeUserHash
}
func (r *MockUserRepository) CreateUser(ctx context.Context, userName string) (*model.User, error) {
	r.CallCreateUserCt++
	return r.RetCreateUser, r.ErrCreateUser
}
func (r *MockUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	r.CallGetByIDCt++
	return r.RetGetById, r.ErrGetById
}
func (r *MockUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	r.CallDeleteUserCt++
	return r.ErrDeleteUser
}

func (r *MockUserRepository) Reset() {
	r.CallMakeUserHashCt = 0
	r.CallCreateUserCt = 0
	r.CallGetByIDCt = 0
	r.CallDeleteUserCt = 0
	r.RetMakeUserHash = ""
	r.RetCreateUser = nil
	r.RetGetById = nil
	r.ErrCreateUser = nil
	r.ErrGetById = nil
	r.ErrDeleteUser = nil
}

// Mock repository for testing
var MockRepository = &MockUserRepository{}

func init() {
	MockRepository = &MockUserRepository{
		CallMakeUserHashCt: 0,
		CallCreateUserCt:   0,
		CallGetByIDCt:      0,
		CallDeleteUserCt:   0,
		RetMakeUserHash:    "",
		RetCreateUser:      nil,
		RetGetById:         nil,
		ErrCreateUser:      nil,
		ErrGetById:         nil,
		ErrDeleteUser:      nil,
	}
}
