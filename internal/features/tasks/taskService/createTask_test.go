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

func TestTaskService_CreateTask(t *testing.T) {
	type args struct {
		ctx        context.Context
		taskDomain domain.Task
	}

	title := "Стать сеньором за 10 минут"
	inputTask := domain.Task{Title: title}
	wantTask := domain.Task{Id: 1, Title: title}
	repoErr := core_errors.TestFailedRepositoryUnknowError

	tests := []struct {
		name         string
		args         args
		want         domain.Task
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockTaskRepository, args args)
	}{
		{
			name: "Успешное создание задачи",
			args: args{
				ctx:        context.Background(),
				taskDomain: inputTask,
			},
			want:    wantTask,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					CreateTask(args.ctx, args.taskDomain).
					Return(wantTask, nil)
			},
		},
		{
			name: "Ошибка создание задачи",
			args: args{
				ctx:        context.Background(),
				taskDomain: inputTask,
			},
			want:    domain.Task{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					CreateTask(args.ctx, args.taskDomain).
					Return(domain.Task{}, repoErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			mockTx := mocks.NewMockTxManager(ctrl)
			tt.mockBehavior(mockRepo, tt.args)

			s := NewTaskService(mockRepo, mockTx)

			got, err := s.CreateTask(tt.args.ctx, tt.args.taskDomain)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTask() got = %v, want %v", got, tt.want)
			}
		})
	}
}
