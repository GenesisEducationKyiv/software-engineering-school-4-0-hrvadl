//go:build !integration

package cron

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/cron/mocks"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

func TestNewAdapter(t *testing.T) {
	t.Parallel()
	type args struct {
		rg      RateGetter
		sg      SubscribersGetter
		s       Sender
		f       Formatter
		timeout time.Duration
		log     *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *Adapter
	}{
		{
			name: "Should create new adapter correctly",
			args: args{
				rg:      mocks.NewMockRateGetter(gomock.NewController(t)),
				sg:      mocks.NewMockSubscribersGetter(gomock.NewController(t)),
				s:       mocks.NewMockSender(gomock.NewController(t)),
				f:       mocks.NewMockFormatter(gomock.NewController(t)),
				timeout: time.Second * 5,
				log:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			},
			want: &Adapter{
				rate:        mocks.NewMockRateGetter(gomock.NewController(t)),
				subscribers: mocks.NewMockSubscribersGetter(gomock.NewController(t)),
				sender:      mocks.NewMockSender(gomock.NewController(t)),
				formatter:   mocks.NewMockFormatter(gomock.NewController(t)),
				timeout:     time.Second * 5,
				log:         slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewAdapter(
				tt.args.rg,
				tt.args.sg,
				tt.args.s,
				tt.args.f,
				tt.args.timeout,
				tt.args.log,
			)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAdapterDo(t *testing.T) {
	t.Parallel()
	type fields struct {
		rate        func(*gomock.Controller) RateGetter
		subscribers func(*gomock.Controller) SubscribersGetter
		sender      func(*gomock.Controller) Sender
		formatter   func(*gomock.Controller) Formatter
		timeout     time.Duration
		log         *slog.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Should return error when rate getter failed",
			fields: fields{
				rate: func(c *gomock.Controller) RateGetter {
					rg := mocks.NewMockRateGetter(c)
					rg.EXPECT().
						Get(gomock.Any()).
						Times(1).
						Return(nil, errors.New("failed to get rate"))
					return rg
				},
				subscribers: func(c *gomock.Controller) SubscribersGetter {
					sg := mocks.NewMockSubscribersGetter(c)
					sg.EXPECT().GetAll(gomock.Any()).Times(0)
					return sg
				},
				sender: func(c *gomock.Controller) Sender {
					s := mocks.NewMockSender(c)
					s.EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
					return s
				},
				formatter: func(c *gomock.Controller) Formatter {
					f := mocks.NewMockFormatter(c)
					f.EXPECT().Format(gomock.Any()).Times(0)
					return f
				},
				log:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				timeout: time.Second * 1,
			},
			wantErr: true,
		},
		{
			name: "Should return error when sub getter failed",
			fields: fields{
				rate: func(c *gomock.Controller) RateGetter {
					rg := mocks.NewMockRateGetter(c)
					rg.EXPECT().
						Get(gomock.Any()).
						Times(1).
						Return(&rate.Exchange{Rate: 32.2}, nil)
					return rg
				},
				subscribers: func(c *gomock.Controller) SubscribersGetter {
					sg := mocks.NewMockSubscribersGetter(c)
					sg.EXPECT().GetAll(gomock.Any()).Times(1).
						Return(nil, errors.New("failed to get subs"))
					return sg
				},
				sender: func(c *gomock.Controller) Sender {
					s := mocks.NewMockSender(c)
					s.EXPECT().
						Send(gomock.Any(), gomock.Any()).
						Times(0)
					return s
				},
				formatter: func(c *gomock.Controller) Formatter {
					f := mocks.NewMockFormatter(c)
					f.EXPECT().Format(gomock.Any()).Times(0)
					return f
				},
				log:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				timeout: time.Second * 1,
			},
			wantErr: true,
		},
		{
			name: "Should return error when sender failed",
			fields: fields{
				rate: func(c *gomock.Controller) RateGetter {
					rg := mocks.NewMockRateGetter(c)
					rg.EXPECT().
						Get(gomock.Any()).
						Times(1).
						Return(&rate.Exchange{Rate: 32.2}, nil)
					return rg
				},
				subscribers: func(c *gomock.Controller) SubscribersGetter {
					sg := mocks.NewMockSubscribersGetter(c)
					sg.EXPECT().GetAll(gomock.Any()).Times(1).
						Return([]subscriber.Subscriber{{Email: "test@test.com"}}, nil)
					return sg
				},
				sender: func(c *gomock.Controller) Sender {
					s := mocks.NewMockSender(c)
					s.EXPECT().
						Send(gomock.Any(), "rate 32.2").
						Times(1).
						Return(errors.New("failed to send"))
					return s
				},
				formatter: func(c *gomock.Controller) Formatter {
					f := mocks.NewMockFormatter(c)
					f.EXPECT().Format(float32(32.2)).Times(1).Return("rate 32.2")
					return f
				},
				log:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				timeout: time.Second * 1,
			},
			wantErr: true,
		},
		{
			name: "Should not return error ",
			fields: fields{
				rate: func(c *gomock.Controller) RateGetter {
					rg := mocks.NewMockRateGetter(c)
					rg.EXPECT().
						Get(gomock.Any()).
						Times(1).
						Return(&rate.Exchange{Rate: 32.2}, nil)
					return rg
				},
				subscribers: func(c *gomock.Controller) SubscribersGetter {
					sg := mocks.NewMockSubscribersGetter(c)
					sg.EXPECT().GetAll(gomock.Any()).Times(1).
						Return([]subscriber.Subscriber{{Email: "test@test.com"}}, nil)
					return sg
				},
				sender: func(c *gomock.Controller) Sender {
					s := mocks.NewMockSender(c)
					s.EXPECT().
						Send(gomock.Any(), "rate 32.2").
						Times(1).
						Return(nil)
					return s
				},
				formatter: func(c *gomock.Controller) Formatter {
					f := mocks.NewMockFormatter(c)
					f.EXPECT().Format(float32(32.2)).Times(1).Return("rate 32.2")
					return f
				},
				log:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				timeout: time.Second * 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := gomock.NewController(t)
			a := &Adapter{
				rate:        tt.fields.rate(c),
				subscribers: tt.fields.subscribers(c),
				sender:      tt.fields.sender(c),
				formatter:   tt.fields.formatter(c),
				timeout:     tt.fields.timeout,
				log:         tt.fields.log,
			}

			err := a.Do()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
