package archtest

import (
	"testing"

	archgo "github.com/fdaines/arch-go/api"
	config "github.com/fdaines/arch-go/api/configuration"
	"github.com/stretchr/testify/require"
)

const moduleInfo = "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub"

func TestDependencies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		cfg  config.Config
	}{
		{
			name: "cmd should depend only on internal/cmd/pkg",
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
			name: "internal should not depend cmd",
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
			name: "pkg should not depend cmd/internal",
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
			name: "platform should not depend on other pkg",
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
		{
			name: "storage should depend only on platform",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.storage.**",
						ShouldOnlyDependsOn: &config.Dependencies{
							Internal: []string{"**.platform.**"},
						},
						ShouldNotDependsOn: &config.Dependencies{
							Internal: []string{"**.transport.**", "**.app.**", "**.service.**"},
						},
					},
				},
			},
		},
		{
			name: "service should depend only on service & storage & pkg",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.service.**",
						ShouldOnlyDependsOn: &config.Dependencies{
							Internal: []string{"**.service.**", "**.storage.**", "**.pkg.**"},
						},
						ShouldNotDependsOn: &config.Dependencies{
							Internal: []string{"**.transport.**", "**.app.**"},
						},
					},
				},
			},
		},
		{
			name: "transport should depend only on service",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.transport.**",
						ShouldOnlyDependsOn: &config.Dependencies{
							Internal: []string{"**.service.**"},
						},
						ShouldNotDependsOn: &config.Dependencies{
							Internal: []string{
								"**.transport.**",
								"**.app.**",
								"**.storage.**",
								"**.pkg.**",
							},
						},
					},
				},
			},
		},
		{
			name: "config should be independent",
			cfg: config.Config{
				DependenciesRules: []*config.DependenciesRule{
					{
						Package: "**.cfg.**",
						ShouldNotDependsOn: &config.Dependencies{
							Internal: []string{
								"**.transport.**",
								"**.app.**",
								"**.service.**",
								"**.storage.**",
								"**.pkg.**",
							},
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
