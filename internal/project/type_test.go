package project

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
)

func TestProjectConfig_AddonGroups(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]Environment
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "should return a list of addon groups",
			fields: fields{
				Addons: map[string]Addon{
					"addon1": {
						Group: "group1",
					},
					"addon2": {
						Group: "group2",
					},
					"addon3": {
						Group: "group1",
					},
				},
			},
			want: []string{"group1", "group2"},
		},
		{
			name: "should return an empty list if no addons are defined",
			fields: fields{
				Addons: map[string]Addon{},
			},
			want: []string{},
		},
		{
			name: "should return a list of addon groups if only one group is defined",
			fields: fields{
				Addons: map[string]Addon{
					"addon1": {
						Group: "group1",
					},
					"addon2": {
						Group: "group1",
					},
				},
			},
			want: []string{"group1"},
		},
		{
			name: "should return an empty list no group is defined",
			fields: fields{
				Addons: map[string]Addon{
					"addon1": {},
					"addon2": {},
				},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProjectConfig{
				BasePath:         tt.fields.BasePath,
				TemplateBasePath: tt.fields.TemplateBasePath,
				Addons:           tt.fields.Addons,
				ParsedAddons:     tt.fields.ParsedAddons,
				Environments:     tt.fields.Environments,
			}
			got := p.AddonGroups()
			slices.Sort(got)
			slices.Sort(tt.want)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ProjectConfig.AddonGroups() mismatch (-got +want):\n%s", diff)
				return
			}
		})
	}
}

func TestProjectConfig_EnvStageProperty(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]Environment
	}
	type args struct {
		environment string
		stage       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "should merge the properties of the environment and stage",
			fields: fields{
				Environments: map[string]Environment{
					"env1": {
						Properties: map[string]string{
							"env1": "env1",
						},
						Stages: map[string]Stage{
							"stage1": {
								Properties: map[string]string{
									"stage1": "stage1",
								},
							},
						},
					},
				},
			},
			args: args{
				environment: "env1",
				stage:       "stage1",
			},
			want: map[string]string{
				"env1":   "env1",
				"stage1": "stage1",
			},
		},
		{
			name: "should return the properties of the stage as the stage has precedence",
			fields: fields{
				Environments: map[string]Environment{
					"env1": {
						Properties: map[string]string{
							"env1": "env1",
						},
					},
				},
			},
			args: args{
				environment: "env1",
				stage:       "stage1",
			},
			want: map[string]string{
				"env1": "env1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &ProjectConfig{
				BasePath:         tt.fields.BasePath,
				TemplateBasePath: tt.fields.TemplateBasePath,
				Addons:           tt.fields.Addons,
				ParsedAddons:     tt.fields.ParsedAddons,
				Environments:     tt.fields.Environments,
			}

			got := pc.EnvStageProperty(tt.args.environment, tt.args.stage)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ProjectConfig.EnvStageProperty() mismatch (-got +want):\n%s", diff)
				return
			}
		})
	}
}

func TestProjectConfig_EnvStageClusterProperty(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]Environment
	}
	type args struct {
		environment string
		stage       string
		cluster     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "should merge the properties of the environment, stage and cluster",
			fields: fields{
				Environments: map[string]Environment{
					"env1": {
						Properties: map[string]string{
							"env1": "env1",
						},
						Stages: map[string]Stage{
							"stage1": {
								Properties: map[string]string{
									"stage1": "stage1",
								},
								Clusters: map[string]Cluster{
									"cluster1": {
										Properties: map[string]string{
											"cluster1": "cluster1",
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				environment: "env1",
				stage:       "stage1",
				cluster:     "cluster1",
			},
			want: map[string]string{
				"env1":     "env1",
				"stage1":   "stage1",
				"cluster1": "cluster1",
			},
		},
		{
			name: "should use the properties of the cluster as the cluster has precedence",
			fields: fields{
				Environments: map[string]Environment{
					"env1": {
						Properties: map[string]string{
							"env1": "env1",
						},
						Stages: map[string]Stage{
							"stage1": {
								Properties: map[string]string{
									"stage1": "stage1",
								},
								Clusters: map[string]Cluster{
									"cluster1": {
										Properties: map[string]string{
											"env1":     "b",
											"stage1":   "a",
											"cluster1": "c",
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				environment: "env1",
				stage:       "stage1",
				cluster:     "cluster1",
			},
			want: map[string]string{
				"env1":     "b",
				"stage1":   "a",
				"cluster1": "c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &ProjectConfig{
				BasePath:         tt.fields.BasePath,
				TemplateBasePath: tt.fields.TemplateBasePath,
				Addons:           tt.fields.Addons,
				ParsedAddons:     tt.fields.ParsedAddons,
				Environments:     tt.fields.Environments,
			}
			got := pc.EnvStageClusterProperty(tt.args.environment, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ProjectConfig.EnvStageClusterProperty() mismatch (-got +want):\n%s", diff)
				return
			}
		})
	}
}
