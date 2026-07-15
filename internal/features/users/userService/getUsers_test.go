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

func TestUserService_GetUsers(t *testing.T) {
	type args struct {
		ctx    context.Context
		limit  *int
		offset *int
	}

	limit := 10
	offset := 0
	users := []domain.User{
		domain.User{Id: 9},
		domain.User{Id: 10},
	}
	repoErr := core_errors.TestFailedRepositoryUnknowError

	tests := []struct {
		name         string
		args         args
		want         []domain.User
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockUserRepository, args args)
	}{
		{
			name: "Успешное получение пользователей с паггинацией",
			args: args{
				ctx:    context.Background(),
				limit:  &limit,
				offset: &offset,
			},
			want:    users,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					GetUsers(args.ctx, args.limit, args.offset).
					Return(users, nil)
			},
		},
		{
			name: "Не получилось получить пользователей с паггинацией",
			args: args{
				ctx:    context.Background(),
				limit:  &limit,
				offset: &offset,
			},
			want:    nil,
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					GetUsers(args.ctx, args.limit, args.offset).
					Return(nil, repoErr)
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

			got, err := s.GetUsers(tt.args.ctx, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
