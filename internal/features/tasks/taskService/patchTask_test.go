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

func TestTaskService_PatchTask(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    int
		patch domain.TaskPatch
	}
	inputId := 1
	oldDescription := "old description"
	title := "test"

	patch := domain.NewTaskPatch(
		domain.Nullable[string]{Value: &title, Set: true},
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[bool]{Value: nil, Set: false},
	)

	existingTask := domain.Task{
		Id:          inputId,
		Version:     1,
		Title:       "old title",
		Description: &oldDescription,
		Completed:   false,
		UserId:      1,
	}

	mergedTask := domain.Task{
		Id:          inputId,
		Version:     1,
		Title:       "test",
		Description: &oldDescription,
		Completed:   false,
		UserId:      1,
	}

	wantTask := domain.Task{
		Id:          inputId,
		Version:     2,
		Title:       "test",
		Description: &oldDescription,
		Completed:   false,
		UserId:      1,
	}

	errNotFound := core_errors.ErrNotFound
	errConflict := core_errors.ErrConflict

	tests := []struct {
		name         string
		args         args
		want         domain.Task
		wantErr      bool
		mockBehavior func(mockRepo *mocks.MockTaskRepository, args args)
	}{
		{
			name: "Успешно обновление задачи",
			args: args{
				ctx:   context.Background(),
				id:    inputId,
				patch: patch,
			},
			want:    wantTask,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					GetTask(args.ctx, args.id).
					Return(existingTask, nil)

				mockRepo.EXPECT().
					PatchTask(args.ctx, args.id, mergedTask).
					Return(wantTask, nil)
			},
		},
		{
			name: "Не удалось обновить задачу(ErrNotFound)",
			args: args{
				ctx:   context.Background(),
				id:    inputId,
				patch: patch,
			},
			want:    domain.Task{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					GetTask(args.ctx, args.id).
					Return(domain.Task{}, errNotFound)
			},
		},
		{
			name: "Не удалось обновить задачу(ErrConflict)",
			args: args{
				ctx:   context.Background(),
				id:    inputId,
				patch: patch,
			},
			want:    domain.Task{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockTaskRepository, args args) {
				mockRepo.EXPECT().
					GetTask(args.ctx, args.id).
					Return(existingTask, nil)

				mockRepo.EXPECT().
					PatchTask(args.ctx, args.id, mergedTask).
					Return(domain.Task{}, errConflict)
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

			got, err := s.PatchTask(tt.args.ctx, tt.args.id, tt.args.patch)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchTask() got = %v, want %v", got, tt.want)
			}
		})
	}
}
