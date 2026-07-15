package statService

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/features/statistics/statService/mocks"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestStatisticsService_GetStatistics(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId *int
		from   *time.Time
		to     *time.Time
	}

	userId := 1
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	created1 := time.Date(2026, 1, 5, 10, 0, 0, 0, time.UTC)
	completed1 := created1.Add(1 * time.Hour)

	created2 := time.Date(2026, 1, 6, 10, 0, 0, 0, time.UTC)
	completed2 := created2.Add(3 * time.Hour)

	created3 := time.Date(2026, 1, 7, 10, 0, 0, 0, time.UTC) // задача 3: не выполнена

	tasks := []domain.Task{
		{Id: 1, Completed: true, CreatedAt: created1, CompletedAt: &completed1},
		{Id: 2, Completed: true, CreatedAt: created2, CompletedAt: &completed2},
		{Id: 3, Completed: false, CreatedAt: created3, CompletedAt: nil},
	}

	wantRate := float64(2) / float64(3) * 100
	wantAvg := 2 * time.Hour
	wantStatistics := domain.Statistics{
		TaskCreated:              3,
		TaskCompleted:            2,
		TaskCompletedRate:        &wantRate,
		TaskAverageCompletedTime: &wantAvg,
	}

	repoErr := errors.New("unexpected repository error")

	tests := []struct {
		name         string
		args         args
		want         domain.Statistics
		wantErr      bool
		wantErrIs    error
		mockBehavior func(mockRepo *mocks.MockStatisticsRepository, args args)
	}{
		{
			name: "Успешное получение статистики",
			args: args{
				ctx:    context.Background(),
				userId: &userId,
				from:   &from,
				to:     &to,
			},
			want:    wantStatistics,
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockStatisticsRepository, args args) {
				mockRepo.EXPECT().
					GetTasks(args.ctx, args.userId, args.from, args.to).
					Return(tasks, nil)
			},
		},
		{
			name: "Пустой список задач",
			args: args{
				ctx:    context.Background(),
				userId: &userId,
				from:   &from,
				to:     &to,
			},
			want:    domain.Statistics{},
			wantErr: false,
			mockBehavior: func(mockRepo *mocks.MockStatisticsRepository, args args) {
				mockRepo.EXPECT().
					GetTasks(args.ctx, args.userId, args.from, args.to).
					Return([]domain.Task{}, nil)
			},
		},
		{
			name: "Ошибка валидации: to раньше from",
			args: args{
				ctx:    context.Background(),
				userId: &userId,
				from:   &to,
				to:     &from,
			},
			want:      domain.Statistics{},
			wantErr:   true,
			wantErrIs: core_errors.ErrInvalidArgument,
			mockBehavior: func(mockRepo *mocks.MockStatisticsRepository, args args) {
			},
		},
		{
			name: "Ошибка валидации: to равен from",
			args: args{
				ctx:    context.Background(),
				userId: &userId,
				from:   &from,
				to:     &from,
			},
			want:      domain.Statistics{},
			wantErr:   true,
			wantErrIs: core_errors.ErrInvalidArgument,
			mockBehavior: func(mockRepo *mocks.MockStatisticsRepository, args args) {
			},
		},
		{
			name: "Ошибка репозитория",
			args: args{
				ctx:    context.Background(),
				userId: &userId,
				from:   &from,
				to:     &to,
			},
			want:    domain.Statistics{},
			wantErr: true,
			mockBehavior: func(mockRepo *mocks.MockStatisticsRepository, args args) {
				mockRepo.EXPECT().
					GetTasks(args.ctx, args.userId, args.from, args.to).
					Return(nil, repoErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockStatisticsRepository(ctrl)
			tt.mockBehavior(mockRepo, tt.args)

			s := NewStatisticsService(mockRepo)

			got, err := s.GetStatistics(tt.args.ctx, tt.args.userId, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatistics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("GetStatistics() error = %v, want error Is %v", err, tt.wantErrIs)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStatistics() got = %v, want %v", got, tt.want)
			}
		})
	}
}
