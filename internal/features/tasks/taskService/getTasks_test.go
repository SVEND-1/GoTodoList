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

func TestTaskService_GetTasks(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId *int
		limit  *int
		offset *int
	}
	limit := 10
	offset := 0
	tasks := []domain.Task{
		domain.Task{Id: 1},
		domain.Task{Id: 2},
	}
	repoErr := core_errors.TestFailedRepositoryUnknowError

	tests := []struct {
		name         string
		args         args
		want         []domain.Task
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockTaskRepository, args args)
	}{
		{
			name: "Успешное получение задач с паггинацией",
			args: args{
				ctx:    context.Background(),
				limit:  &limit,
				offset: &offset,
			},
			want:    tasks,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					GetTasks(args.ctx, args.userId, args.limit, args.offset).
					Return(tasks, nil)
			},
		},
		{
			name: "Не получилось получить задач с паггинацией",
			args: args{
				ctx:    context.Background(),
				limit:  &limit,
				offset: &offset,
			},
			want:    nil,
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					GetTasks(args.ctx, args.userId, args.limit, args.offset).
					Return(nil, repoErr)
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

			got, err := s.GetTasks(tt.args.ctx, tt.args.userId, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTasks() got = %v, want %v", got, tt.want)
			}
		})
	}
}
