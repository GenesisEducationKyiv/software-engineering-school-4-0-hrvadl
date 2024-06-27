package archtest

import (
	"testing"

	archgo "github.com/fdaines/arch-go/api"
	config "github.com/fdaines/arch-go/api/configuration"
	"github.com/stretchr/testify/require"
)

const moduleInfo = "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer"

func TestDependencies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		cfg  config.Config
	}{
		{
			name: "cmd should depends only on internal/cmd/pkg",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.cmd.**",
						ShouldOnlyDependsOn: &config.Dependencies{
							Internal: []string{"**.cmd.**", "**.internal.**", "**.pkg.*"},
						},
					},
				},
			},
		},
		{
			name: "internal should not depends cmd",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.internal.**",
						ShouldNotDependsOn: &config.Dependencies{
							Internal: []string{"**.cmd.**"},
						},
					},
				},
			},
		},
		{
			name: "pkg should not depends cmd/internal",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.pkg.**",
						ShouldNotDependsOn: &config.Dependencies{
							Internal: []string{"**.cmd.**", "**.internal.**"},
						},
					},
				},
			},
		},
		{
			name: "platform should not depends on other pkg",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.platform.**",
						ShouldNotDependsOn: &config.Dependencies{
							Internal: []string{"**.transport.**", "**.app.**", "**.platform.**"},
						},
					},
				},
			},
		},
	}

	module := config.Load(moduleInfo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.True(t, archgo.CheckArchitecture(module, tt.cfg).Passes)
		})
	}
}
