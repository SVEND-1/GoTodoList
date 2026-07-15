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

func TestUserService_GetUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	inputId := 1
	wantUserSuccess := domain.User{Id: inputId}
	repoErr := core_errors.ErrNotFound

	tests := []struct {
		name         string
		args         args
		want         domain.User
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockUserRepository, args args)
	}{
		{
			name: "Пользователь найден успешно",
			args: args{
				ctx: context.Background(),
				id:  inputId,
			},
			want:    wantUserSuccess,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					GetUser(args.ctx, args.id).
					Return(wantUserSuccess, nil)
			},
		},
		{
			name: "Пользователь не найден",
			args: args{
				ctx: context.Background(),
				id:  inputId,
			},
			want:    domain.User{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					GetUser(args.ctx, args.id).
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

			got, err := s.GetUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
