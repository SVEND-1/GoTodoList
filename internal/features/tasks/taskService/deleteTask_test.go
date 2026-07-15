package taskService

import (
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/features/tasks/taskService/mocks"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestTaskService_DeleteTask(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	inputId := 1
	repoErr := core_errors.ErrNotFound

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockTaskRepository, args args)
	}{
		{
			name: "Успешное удаление задачи",
			args: args{
				ctx: context.Background(),
				id:  inputId,
			},
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					DeleteTask(args.ctx, args.id).
					Return(nil)
			},
		},
		{
			name: "Ошибка удаление задачи",
			args: args{
				ctx: context.Background(),
				id:  inputId,
			},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					DeleteTask(args.ctx, args.id).
					Return(repoErr)
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

			if err := s.DeleteTask(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
