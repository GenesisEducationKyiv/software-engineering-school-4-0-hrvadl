package cron

import (
	"log/slog"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/hrvadl/btcratenotifier/sub/internal/service/cron/mocks"
)

func TestNewDailyJob(t *testing.T) {
	t.Parallel()
	type args struct {
		hour int
		min  int
		log  *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *Job
	}{
		{
			name: "Should return correct daily job when arguments are correct",
			args: args{
				hour: 12,
				min:  0,
				log:  slog.Default(),
			},
			want: &Job{
				log:      slog.Default(),
				interval: time.Hour * 24,
			},
		},
		{
			name: "Should return correct daily job when arguments are allowed",
			args: args{
				hour: 0,
				min:  0,
				log:  nil,
			},
			want: &Job{
				log:      nil,
				interval: time.Hour * 24,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewDailyJob(tt.args.hour, tt.args.min, tt.args.log); !reflect.DeepEqual(
				got.interval,
				tt.want.interval,
			) {
				t.Errorf("NewDailyJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJob(t *testing.T) {
	t.Parallel()
	type args struct {
		interval time.Duration
		log      *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *Job
	}{
		{
			name: "Should create job correctly when arguments are correct",
			args: args{
				interval: time.Hour * 2,
				log:      slog.Default(),
			},
			want: &Job{
				interval: time.Hour * 2,
				log:      slog.Default(),
			},
		},
		{
			name: "Should create job correctly when arguments are correct",
			args: args{
				interval: time.Hour * 1,
				log:      nil,
			},
			want: &Job{
				interval: time.Hour * 1,
				log:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewJob(tt.args.interval, tt.args.log); !reflect.DeepEqual(
				got.interval,
				tt.want.interval,
			) {
				t.Errorf("NewJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobDo(t *testing.T) {
	t.Parallel()
	type fields struct {
		interval time.Duration
		log      *slog.Logger
	}
	type args struct {
		doer Doer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, doer Doer)
	}{
		{
			name: "Should execute job correcty with 5ms timeout",
			fields: fields{
				interval: time.Millisecond * 3,
				log:      slog.Default(),
			},
			args: args{
				doer: mocks.NewMockDoer(gomock.NewController(t)),
			},
			setup: func(t *testing.T, doer Doer) {
				t.Helper()
				d, ok := doer.(*mocks.MockDoer)
				if !ok {
					t.Fatalf("failed to cast doer to mock doer")
				}
				d.EXPECT().Do().MinTimes(3).Return(nil)
			},
		},
		{
			name: "Should execute job correcty with 3ms timeout",
			fields: fields{
				interval: time.Millisecond,
				log:      slog.Default(),
			},
			args: args{
				doer: mocks.NewMockDoer(gomock.NewController(t)),
			},
			setup: func(t *testing.T, doer Doer) {
				t.Helper()
				d, ok := doer.(*mocks.MockDoer)
				if !ok {
					t.Fatalf("failed to cast doer to mock doer")
				}
				d.EXPECT().Do().MinTimes(3).Return(nil)
			},
		},
		{
			name: "Should execute job correcty with 3ms timeout",
			fields: fields{
				interval: time.Millisecond * 3,
				log:      slog.Default(),
			},
			args: args{
				doer: mocks.NewMockDoer(gomock.NewController(t)),
			},
			setup: func(t *testing.T, doer Doer) {
				t.Helper()
				d, ok := doer.(*mocks.MockDoer)
				if !ok {
					t.Fatalf("failed to cast doer to mock doer")
				}
				d.EXPECT().Do().MinTimes(3).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.args.doer)
			j := NewJob(tt.fields.interval, tt.fields.log)
			j.Do(tt.args.doer)
			time.Sleep(tt.fields.interval * 10)
		})
	}
}
