package taskService

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/features/tasks/taskService/mocks"
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestTaskService_GetTask(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	inputId := 1
	wantTask := domain.Task{Id: inputId}
	repoErr := core_errors.ErrNotFound

	tests := []struct {
		name         string
		args         args
		want         domain.Task
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockTaskRepository, args args)
	}{
		{
			name: "Задача найдена успешно",
			args: args{
				ctx: context.Background(),
				id:  inputId,
			},
			want:    wantTask,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					GetTask(args.ctx, args.id).
					Return(domain.Task{Id: inputId}, nil)
			},
		},
		{
			name: "Задача не найдена",
			args: args{
				ctx: context.Background(),
				id:  inputId,
			},
			want:    domain.Task{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					GetTask(args.ctx, args.id).
					Return(domain.Task{}, repoErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			tt.mockBehavior(mockRepo, tt.args)

			s := NewTaskService(mockRepo)

			got, err := s.GetTask(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTask() got = %v, want %v", got, tt.want)
			}
		})
	}
}
