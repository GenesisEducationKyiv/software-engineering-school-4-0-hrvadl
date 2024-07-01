package mail

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	model "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/models/mail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/mail/mocks"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	type args struct {
		m Mailer
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Should create service with correct arguments",
			args: args{
				m: mocks.NewMockMailer(gomock.NewController(t)),
			},
			want: &Service{
				mailers: []Mailer{mocks.NewMockMailer(gomock.NewController(t))},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewService(tt.args.m)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceSend(t *testing.T) {
	t.Parallel()
	type fields struct {
		mailers func(*gomock.Controller) []Mailer
	}
	type args struct {
		ctx  context.Context
		mail model.Mail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not fallback when first mailer succeeded",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						To:      []string{"test@test.com"},
						Subject: "Sub",
						HTML:    "<h1>Hello</h1>",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(1).Return(nil)

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(0)

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(0)

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					To:      []string{"test@test.com"},
					Subject: "Sub",
					HTML:    "<h1>Hello</h1>",
				},
			},
			wantErr: false,
		},
		{
			name: "Should fallback when first mailer failed",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						To:      []string{"test@test.com"},
						Subject: "Sub",
						HTML:    "<h1>Hello</h1>",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(1).Return(errors.New("failed"))

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(1).Return(nil)

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(0)

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					To:      []string{"test@test.com"},
					Subject: "Sub",
					HTML:    "<h1>Hello</h1>",
				},
			},
			wantErr: false,
		},
		{
			name: "Should fallback to third, when second failed",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						To:      []string{"test@test.com"},
						Subject: "Sub",
						HTML:    "<h1>Hello</h1>",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(1).Return(errors.New("failed"))

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(1).Return(errors.New("failed"))

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(1).Return(nil)

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					To:      []string{"test@test.com"},
					Subject: "Sub",
					HTML:    "<h1>Hello</h1>",
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error, when all failed",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						To:      []string{"test@test.com"},
						Subject: "Sub",
						HTML:    "<h1>Hello</h1>",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(1).Return(errors.New("failed"))

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(1).Return(errors.New("failed"))

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(1).Return(errors.New("failed"))

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					To:      []string{"test@test.com"},
					Subject: "Sub",
					HTML:    "<h1>Hello</h1>",
				},
			},
			wantErr: true,
		},
		{
			name: "Should not be called, when body is empty",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						To:      []string{"test@test.com"},
						Subject: "Sub",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(0)

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(0)

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(0)

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					To:      []string{"test@test.com"},
					Subject: "Sub",
				},
			},
			wantErr: true,
		},
		{
			name: "Should not be called, when subject is empty",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						To:   []string{"test@test.com"},
						HTML: "<h1>Hello</h1>",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(0)

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(0)

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(0)

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					To:   []string{"test@test.com"},
					HTML: "<h1>Hello</h1>",
				},
			},
			wantErr: true,
		},
		{
			name: "Should not be called, when to field is empty",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						To:      make([]string, 0),
						Subject: "Sub",
						HTML:    "<h1>Hello</h1>",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(0)

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(0)

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(0)

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					To:      make([]string, 0),
					Subject: "Sub",
					HTML:    "<h1>Hello</h1>",
				},
			},
			wantErr: true,
		},
		{
			name: "Should not be called, when to field is nil",
			fields: fields{
				mailers: func(c *gomock.Controller) []Mailer {
					m := model.Mail{
						Subject: "Sub",
						HTML:    "<h1>Hello</h1>",
					}

					m1 := mocks.NewMockMailer(c)
					m1.EXPECT().Send(gomock.Any(), m).Times(0)

					m2 := mocks.NewMockMailer(c)
					m2.EXPECT().Send(gomock.Any(), m).Times(0)

					m3 := mocks.NewMockMailer(c)
					m3.EXPECT().Send(gomock.Any(), m).Times(0)

					return []Mailer{m1, m2, m3}
				},
			},
			args: args{
				mail: model.Mail{
					Subject: "Sub",
					HTML:    "<h1>Hello</h1>",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Service{
				mailers: tt.fields.mailers(gomock.NewController(t)),
			}

			err := s.Send(tt.args.ctx, tt.args.mail)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestServiceSetNext(t *testing.T) {
	t.Parallel()
	type fields struct {
		mailers []Mailer
	}
	type args struct {
		mailers []Mailer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Service
	}{
		{
			name: "Should construct chain of responsibility correctly",
			fields: fields{
				mailers: nil,
			},
			args: args{
				mailers: []Mailer{
					mocks.NewMockMailer(gomock.NewController(t)),
				},
			},
			want: &Service{
				mailers: []Mailer{
					mocks.NewMockMailer(gomock.NewController(t)),
				},
			},
		},
		{
			name: "Should construct chain of responsibility correctly",
			fields: fields{
				mailers: nil,
			},
			args: args{
				mailers: []Mailer{
					mocks.NewMockMailer(gomock.NewController(t)),
					mocks.NewMockMailer(gomock.NewController(t)),
				},
			},
			want: &Service{
				mailers: []Mailer{
					mocks.NewMockMailer(gomock.NewController(t)),
					mocks.NewMockMailer(gomock.NewController(t)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := &Service{
				mailers: tt.fields.mailers,
			}
			got.SetNext(tt.args.mailers...)
			require.Equal(t, tt.want, got)
		})
	}
}
