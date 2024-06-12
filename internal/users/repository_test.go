package users

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestMakeRepository(t *testing.T) {
	type args struct {
		r UserRepository
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MakeRepository allows dependency injection",
			args: args{
				r: MockRepository,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MakeUserRepository(tt.args.r)
			if MockRepository.CallMakeUserHashCt != 0 {
				t.Errorf("MockRepository.CallMakeUserHashCt = %v, want %v", MockRepository.CallMakeUserHashCt, 0)
			}
			Repository.MakeUserHash(uuid.New())
			if MockRepository.CallMakeUserHashCt != 1 {
				t.Errorf("MockRepository.CallMakeUserHashCt = %v, want %v", MockRepository.CallMakeUserHashCt, 1)
			}
			Repository.CreateUser(context.Background(), "testUserName")
			if MockRepository.CallCreateUserCt != 1 {
				t.Errorf("MockRepository.CallCreateUserCt = %v, want %v", MockRepository.CallMakeUserHashCt, 1)
			}
			Repository.GetByID(context.Background(), uuid.New())
			if MockRepository.CallGetByIDCt != 1 {
				t.Errorf("MockRepository.CallGetByIDCt = %v, want %v", MockRepository.CallMakeUserHashCt, 1)
			}
			Repository.DeleteUser(context.Background(), uuid.New())
			if MockRepository.CallDeleteUserCt != 1 {
				t.Errorf("MockRepository.CallDeleteUserCt = %v, want %v", MockRepository.CallMakeUserHashCt, 1)
			}
			MockRepository.Reset()
			if MockRepository.CallMakeUserHashCt != 0 {
				t.Errorf("MockRepository.CallMakeUserHashCt = %v, want %v", MockRepository.CallMakeUserHashCt, 0)
			}
		})
	}
}

func TestUserRepository_MakeUserHash(t *testing.T) {
	testIDString := "81974462-2d4f-4ce0-b12e-be42a9985838"
	testUUID, _ := uuid.Parse(testIDString)
	testWant := "test-prefix:81974462-2d4f-4ce0-b12e-be42a9985838"
	type fields struct {
		client       *redis.Client
		hashPrefix   string
		keyName      string
		keyCreatedAt string
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "MakeUserHash generates a valid string containing user UUID",
			fields: fields{
				client:       nil,
				hashPrefix:   "test-prefix",
				keyName:      "",
				keyCreatedAt: "",
			},
			args: args{
				userID: testUUID,
			},
			want: testWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepositoryImpl{
				client:       tt.fields.client,
				hashPrefix:   tt.fields.hashPrefix,
				keyName:      tt.fields.keyName,
				keyCreatedAt: tt.fields.keyCreatedAt,
			}
			if got := r.MakeUserHash(tt.args.userID); got != tt.want {
				t.Errorf("UserRepository.MakeUserHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
