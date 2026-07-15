package userService

import (
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/features/users/userService/mocks"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserService_DeleteUser(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId int
	}

	inputId := 1
	repoErr := core_errors.ErrNotFound

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockUserRepository, args args)
	}{
		{
			name: "Успешно удаление пользователя",
			args: args{
				ctx:    context.Background(),
				userId: inputId,
			},
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					DeleteUser(args.ctx, args.userId).
					Return(nil)
			},
		},
		{
			name: "Ошибка удаление пользователя, NotFoundError",
			args: args{
				ctx:    context.Background(),
				userId: inputId,
			},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					DeleteUser(args.ctx, args.userId).
					Return(repoErr)
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
			err := s.DeleteUser(tt.args.ctx, tt.args.userId)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
