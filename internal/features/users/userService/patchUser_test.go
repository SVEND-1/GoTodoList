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

func TestUserService_PatchUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    int
		patch domain.UserPatch
	}

	inputId := 1
	oldPhone := "+79990001122"
	name := "test"

	patch := domain.NewUserPatch(
		domain.Nullable[string]{Value: &name, Set: true},
		domain.Nullable[string]{Value: nil, Set: false},
	)

	existingUser := domain.User{
		Id:          inputId,
		Version:     1,
		FullName:    "old name",
		PhoneNumber: &oldPhone,
	}

	mergedUser := domain.User{
		Id:          inputId,
		Version:     1,
		FullName:    "test",
		PhoneNumber: &oldPhone,
	}

	wantUser := domain.User{
		Id:          inputId,
		Version:     2,
		FullName:    "test",
		PhoneNumber: &oldPhone,
	}

	errNotFound := core_errors.ErrNotFound
	errConflict := core_errors.ErrConflict

	tests := []struct {
		name         string
		args         args
		want         domain.User
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockUserRepository, args args)
	}{
		{
			name: "Успешное обновление пользователя",
			args: args{
				ctx:   context.Background(),
				id:    inputId,
				patch: patch,
			},
			want:    wantUser,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					GetUser(args.ctx, args.id).
					Return(existingUser, nil)

				mockRepo.EXPECT().
					PatchUser(args.ctx, args.id, mergedUser).
					Return(wantUser, nil)
			},
		},
		{
			name: "Не удалось обновить пользователя (ErrNotFound)",
			args: args{
				ctx:   context.Background(),
				id:    inputId,
				patch: patch,
			},
			want:    domain.User{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					GetUser(args.ctx, args.id).
					Return(domain.User{}, errNotFound)
			},
		},
		{
			name: "Не удалось обновить пользователя (ErrConflict)",
			args: args{
				ctx:   context.Background(),
				id:    inputId,
				patch: patch,
			},
			want:    domain.User{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockUserRepository, args args) {
				mockRepo.EXPECT().
					GetUser(args.ctx, args.id).
					Return(existingUser, nil)

				mockRepo.EXPECT().
					PatchUser(args.ctx, args.id, mergedUser).
					Return(domain.User{}, errConflict)
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

			got, err := s.PatchUser(tt.args.ctx, tt.args.id, tt.args.patch)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
