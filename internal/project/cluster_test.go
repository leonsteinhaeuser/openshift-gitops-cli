package project

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
)

func TestClusterAddon_AllRequiredPropertiesSet(t *testing.T) {
	type fields struct {
		Enabled    bool
		Properties map[string]any
	}
	type args struct {
		config    *ProjectConfig
		addonName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no properties defined in addon properties",
			fields: fields{
				Properties: map[string]any{},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeBool,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
		{
			name: "one property set",
			fields: fields{
				Properties: map[string]any{
					"property1": true,
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeBool,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: false,
		},
		{
			name: "one property set string, expect bool",
			fields: fields{
				Properties: map[string]any{
					"property1": "invalid",
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeBool,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect bool",
			fields: fields{
				Properties: map[string]any{
					"property1": 10,
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeBool,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
		{
			name: "one property set bool, expect string",
			fields: fields{
				Properties: map[string]any{
					"property1": true,
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeString,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect string",
			fields: fields{
				Properties: map[string]any{
					"property1": 10,
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeString,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
		{
			name: "one property set string, expect string",
			fields: fields{
				Properties: map[string]any{
					"property1": "value",
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeString,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: false,
		},
		{
			name: "one property set int, expect int",
			fields: fields{
				Properties: map[string]any{
					"property1": 10,
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeInt,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: false,
		},
		{
			name: "one property set string, expect int",
			fields: fields{
				Properties: map[string]any{
					"property1": "value",
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeInt,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect int",
			fields: fields{
				Properties: map[string]any{
					"property1": 10,
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeInt,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: false,
		},
		{
			name: "one property set int, expect int, invalid value",
			fields: fields{
				Properties: map[string]any{
					"property1": "invalid",
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeInt,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect int, invalid value",
			fields: fields{
				Properties: map[string]any{
					"property1": "invalid",
				},
			},
			args: args{
				config: &ProjectConfig{
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     nil,
									Type:        template.PropertyTypeInt,
									Description: "property1",
								},
							},
						},
					},
				},
				addonName: "addon1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := &ClusterAddon{
				Enabled:    tt.fields.Enabled,
				Properties: tt.fields.Properties,
			}
			if err := ca.AllRequiredPropertiesSet(tt.args.config, tt.args.addonName); (err != nil) != tt.wantErr {
				t.Errorf("ClusterAddon.AllRequiredPropertiesSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCluster_IsAddonEnabled(t *testing.T) {
	type fields struct {
		Name       string
		Addons     map[string]*ClusterAddon
		Properties map[string]string
	}
	type args struct {
		addon string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "one addon, enabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: true,
		},
		{
			name: "one addon, disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: false,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: false,
		},
		{
			name: "no addon, disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				addon: "addon1",
			},
			want: false,
		},
		{
			name: "two addons, enabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
					"addon2": {
						Enabled: true,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: true,
		},
		{
			name: "two addons, disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: false,
					},
					"addon2": {
						Enabled: false,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: false,
		},
		{
			name: "two addons, one enabled, one disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: false,
					},
					"addon2": {
						Enabled: true,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: false,
		},
		{
			name: "two addons, one enabled, one disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
					"addon2": {
						Enabled: false,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cluster{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			if got := c.IsAddonEnabled(tt.args.addon); got != tt.want {
				t.Errorf("Cluster.IsAddonEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_EnableAddon(t *testing.T) {
	type fields struct {
		Name       string
		Addons     map[string]*ClusterAddon
		Properties map[string]string
	}
	type args struct {
		addon string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Cluster
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
		},
		{
			name: "one addon, enabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
		},
		{
			name: "one addon, disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: false,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
		},
		{
			name: "two addons, one enabled, one disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
					"addon2": {
						Enabled: false,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
					"addon2": {
						Enabled: false,
					},
				},
			},
		},
		{
			name: "two addons, one enabled, one disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: false,
					},
					"addon2": {
						Enabled: true,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
					"addon2": {
						Enabled: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			c.EnableAddon(tt.args.addon)

			diff := cmp.Diff(tt.want, *c)
			if diff != "" {
				t.Errorf("Cluster.EnableAddon() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestCluster_DisableAddon(t *testing.T) {
	type fields struct {
		Name       string
		Addons     map[string]*ClusterAddon
		Properties map[string]string
	}
	type args struct {
		addon string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Cluster
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{},
			},
		},
		{
			name: "one addon, enabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled:    false,
						Properties: map[string]any{},
					},
				},
			},
		},
		{
			name: "one addon, disabled",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: false,
					},
				},
			},
			args: args{
				addon: "addon1",
			},
			want: Cluster{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled:    false,
						Properties: map[string]any{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			c.DisableAddon(tt.args.addon)

			diff := cmp.Diff(tt.want, *c)
			if diff != "" {
				t.Errorf("Cluster.DisableAddon() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
