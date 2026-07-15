package userService

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/features/users/userService/mocks"
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserService_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user domain.User
	}

	inputUser := domain.User{FullName: "Никита"}
	createdUser := domain.User{Id: 1, FullName: "Никита"}
	repoErr := core_errors.TestFailedRepositoryUnknowError //TODO подумать

	tests := []struct {
		name         string
		args         args
		want         domain.User
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockUserRepository, args args)
	}{
		{
			name: "Успешное создание пользователя",
			args: args{
				ctx:  context.Background(),
				user: inputUser,
			},
			want:    createdUser,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					CreateUser(args.ctx, args.user).
					Return(createdUser, nil)
			},
		},
		{
			name: "Ошибка создание пользователя",
			args: args{
				ctx:  context.Background(),
				user: inputUser,
			},
			want:    domain.User{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					CreateUser(args.ctx, args.user).
					Return(domain.User{}, repoErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)
			tt.mockBehavior(mockRepo, tt.args)

			s := NewUserService(mockRepo)

			got, err := s.CreateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
