package project

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
)

func TestProjectConfig_HasCluster(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]*Environment
	}
	type args struct {
		env     string
		stage   string
		cluster string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true if the cluster exists",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {
								Clusters: map[string]*Cluster{
									"cluster1": {},
								},
							},
						},
					},
				},
			},
			args: args{
				env:     "env1",
				stage:   "stage1",
				cluster: "cluster1",
			},
			want: true,
		},
		{
			name: "should return false if the cluster does not exist",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {
								Clusters: map[string]*Cluster{
									"cluster1": {},
								},
							},
						},
					},
				},
			},
			args: args{
				env:     "env1",
				stage:   "stage1",
				cluster: "cluster2",
			},
			want: false,
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
			if got := p.HasCluster(tt.args.env, tt.args.stage, tt.args.cluster); got != tt.want {
				t.Errorf("ProjectConfig.HasCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectConfig_Cluster(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]*Environment
	}
	type args struct {
		env     string
		stage   string
		cluster string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Cluster
	}{
		{
			name: "should return true if the cluster exists",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {
								Clusters: map[string]*Cluster{
									"cluster1": {
										Name: "cluster1",
										Properties: map[string]string{
											"key": "value",
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				env:     "env1",
				stage:   "stage1",
				cluster: "cluster1",
			},
			want: &Cluster{
				Name: "cluster1",
				Properties: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name: "should return false if the cluster does not exist",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {
								Clusters: map[string]*Cluster{
									"cluster1": {},
								},
							},
						},
					},
				},
			},
			args: args{
				env:     "env1",
				stage:   "stage1",
				cluster: "cluster2",
			},
			want: nil,
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

			got := p.Cluster(tt.args.env, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ProjectConfig.Cluster() mismatch (-got +want):\n%s", diff)
				return
			}
		})
	}
}

func TestProjectConfig_SetCluster(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]*Environment
	}
	type args struct {
		env     string
		stage   string
		cluster *Cluster
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Cluster
	}{
		{
			name: "should set the cluster for the given environment and stage",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {},
						},
					},
				},
			},
			args: args{
				env:   "env1",
				stage: "stage1",
				cluster: &Cluster{
					Name: "cluster1",
					Properties: map[string]string{
						"key": "value",
					},
				},
			},
			want: &Cluster{
				Name: "cluster1",
				Properties: map[string]string{
					"key": "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProjectConfig{
				BasePath:         tt.fields.BasePath,
				TemplateBasePath: tt.fields.TemplateBasePath,
				Addons:           tt.fields.Addons,
				ParsedAddons:     tt.fields.ParsedAddons,
				Environments:     tt.fields.Environments,
			}
			p.SetCluster(tt.args.env, tt.args.stage, tt.args.cluster)

			diff := cmp.Diff(p.Environments[tt.args.env].Stages[tt.args.stage].Clusters[tt.args.cluster.Name], tt.want)
			if diff != "" {
				t.Errorf("ProjectConfig.SetCluster() mismatch (-got +want):\n%s", diff)
				return
			}
		})
	}
}

func TestProjectConfig_DeleteCluster(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]*Environment
	}
	type args struct {
		env     string
		stage   string
		cluster string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*Cluster
	}{
		{
			name: "should delete the cluster for the given environment and stage",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {
								Clusters: map[string]*Cluster{
									"cluster1": {},
								},
							},
						},
					},
				},
			},
			args: args{
				env:     "env1",
				stage:   "stage1",
				cluster: "cluster1",
			},
			want: map[string]*Cluster{},
		},
		{
			name: "should not delete the cluster if it does not exist",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {
								Clusters: map[string]*Cluster{
									"cluster1": {},
								},
							},
						},
					},
				},
			},
			args: args{
				env:     "env1",
				stage:   "stage1",
				cluster: "cluster2",
			},
			want: map[string]*Cluster{
				"cluster1": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProjectConfig{
				BasePath:         tt.fields.BasePath,
				TemplateBasePath: tt.fields.TemplateBasePath,
				Addons:           tt.fields.Addons,
				ParsedAddons:     tt.fields.ParsedAddons,
				Environments:     tt.fields.Environments,
			}
			p.DeleteCluster(tt.args.env, tt.args.stage, tt.args.cluster)

			diff := cmp.Diff(p.Environments[tt.args.env].Stages[tt.args.stage].Clusters, tt.want)
			if diff != "" {
				t.Errorf("ProjectConfig.DeleteCluster() mismatch (-got +want):\n%s", diff)
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
		Environments     map[string]*Environment
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
				Environments: map[string]*Environment{
					"env1": {
						Properties: map[string]string{
							"key1": "value1",
						},
						Stages: map[string]*Stage{
							"stage1": {
								Properties: map[string]string{
									"key2": "value2",
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
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "should return an empty map if no properties are defined",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {},
						},
					},
				},
			},
			args: args{
				environment: "env1",
				stage:       "stage1",
			},
			want: map[string]string{},
		},
		{
			name: "should return the properties of the environment if no stage properties are defined",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Properties: map[string]string{
							"key1": "value1",
						},
						Stages: map[string]*Stage{
							"stage1": {},
						},
					},
				},
			},
			args: args{
				environment: "env1",
				stage:       "stage1",
			},
			want: map[string]string{
				"key1": "value1",
			},
		},
		{
			name: "should return the properties of the stage if no environment properties are defined",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {
								Properties: map[string]string{
									"key2": "value2",
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
				"key2": "value2",
			},
		},
		{
			name: "should return an empty map if no environment and stage properties are defined",
			fields: fields{
				Environments: map[string]*Environment{
					"env1": {
						Stages: map[string]*Stage{
							"stage1": {},
						},
					},
				},
			},
			args: args{
				environment: "env1",
				stage:       "stage1",
			},
			want: map[string]string{},
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

func TestProjectConfig_AddonGroups(t *testing.T) {
	type fields struct {
		BasePath         string
		TemplateBasePath string
		Addons           map[string]Addon
		ParsedAddons     map[string]template.TemplateManifest
		Environments     map[string]*Environment
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
