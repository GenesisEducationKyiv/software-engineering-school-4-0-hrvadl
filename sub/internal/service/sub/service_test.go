//go:build !integration

package sub

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sub/mocks"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	type args struct {
		rr RecipientSource
		vv Validator
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Should create new service correctly when correct arguments are provided",
			args: args{
				rr: mocks.NewMockRecipientSource(gomock.NewController(t)),
				vv: mocks.NewMockValidator(gomock.NewController(t)),
			},
			want: &Service{
				repo:      mocks.NewMockRecipientSource(gomock.NewController(t)),
				validator: mocks.NewMockValidator(gomock.NewController(t)),
			},
		},
		{
			name: "Should create new service correctly when allowed arguments are provided",
			args: args{
				rr: nil,
				vv: nil,
			},
			want: &Service{
				repo: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewService(tt.args.rr, tt.args.vv)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceSubscribe(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo      func(ctrl *gomock.Controller) RecipientSource
		validator func(ctrl *gomock.Controller) Validator
	}
	type args struct {
		ctx  context.Context
		mail string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Should not return err when everything is correct",
			fields: fields{
				validator: func(ctrl *gomock.Controller) Validator {
					v := mocks.NewMockValidator(ctrl)
					v.EXPECT().Validate("mail@gmail.com").Times(1).Return(true)
					return v
				},
				repo: func(ctrl *gomock.Controller) RecipientSource {
					rs := mocks.NewMockRecipientSource(ctrl)
					rs.EXPECT().
						Save(gomock.Any(), subscriber.Subscriber{Email: "mail@gmail.com"}).
						Times(1).
						Return(int64(1), nil)
					return rs
				},
			},
			args: args{
				ctx:  context.Background(),
				mail: "mail@gmail.com",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Should return err when saver returned err",
			fields: fields{
				repo: func(ctrl *gomock.Controller) RecipientSource {
					rs := mocks.NewMockRecipientSource(ctrl)
					rs.EXPECT().
						Save(gomock.Any(), subscriber.Subscriber{Email: "mail@gmail.com"}).
						Times(1).
						Return(int64(0), errors.New("failed to save subscriber"))
					return rs
				},
				validator: func(ctrl *gomock.Controller) Validator {
					v := mocks.NewMockValidator(ctrl)
					v.EXPECT().Validate("mail@gmail.com").Times(1).Return(true)
					return v
				},
			},
			args: args{
				ctx:  context.Background(),
				mail: "mail@gmail.com",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Should return err when validator returned false",
			fields: fields{
				repo: func(ctrl *gomock.Controller) RecipientSource {
					rs := mocks.NewMockRecipientSource(ctrl)
					rs.EXPECT().
						Save(gomock.Any(), subscriber.Subscriber{Email: "mail@gmail.com"}).
						Times(0).
						Return(int64(0), errors.New("failed to save subscriber"))
					return rs
				},
				validator: func(ctrl *gomock.Controller) Validator {
					v := mocks.NewMockValidator(ctrl)
					v.EXPECT().Validate("").Times(1).Return(false)
					return v
				},
			},
			args: args{
				ctx:  context.Background(),
				mail: "",
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			s := &Service{repo: tt.fields.repo(ctrl), validator: tt.fields.validator(ctrl)}
			got, err := s.Subscribe(tt.args.ctx, tt.args.mail)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
