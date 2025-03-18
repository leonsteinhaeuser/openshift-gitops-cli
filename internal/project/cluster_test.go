package project

import (
	"reflect"
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
						Enabled:    true,
						Properties: map[string]any{},
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
						Properties: nil,
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

func TestCluster_GetAddons(t *testing.T) {
	type fields struct {
		Name       string
		Addons     map[string]*ClusterAddon
		Properties map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   ClusterAddons
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			want: map[string]*ClusterAddon{},
		},
		{
			name: "one addon",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
			want: map[string]*ClusterAddon{
				"addon1": {
					Enabled: true,
				},
			},
		},
		{
			name: "two addons",
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
			want: map[string]*ClusterAddon{
				"addon1": {
					Enabled: true,
				},
				"addon2": {
					Enabled: false,
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
			if got := c.GetAddons(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cluster.GetAddons() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_GetAddon(t *testing.T) {
	type fields struct {
		Name       string
		Addons     map[string]*ClusterAddon
		Properties map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ClusterAddon
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				name: "addon1",
			},
			want: nil,
		},
		{
			name: "one addon",
			fields: fields{
				Addons: map[string]*ClusterAddon{
					"addon1": {
						Enabled: true,
					},
				},
			},
			args: args{
				name: "addon1",
			},
			want: &ClusterAddon{
				Enabled: true,
			},
		},
		{
			name: "two addons",
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
				name: "addon1",
			},
			want: &ClusterAddon{
				Enabled: true,
			},
		},
		{
			name: "two addons no found",
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
				name: "addon5",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			if got := c.GetAddon(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cluster.GetAddon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_Render(t *testing.T) {
	type fields struct {
		Name       string
		Addons     map[string]*ClusterAddon
		Properties map[string]string
	}
	type args struct {
		config *ProjectConfig
		env    string
		stage  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			if err := c.Render(tt.args.config, tt.args.env, tt.args.stage); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.Render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCluster_SetDefaultAddons(t *testing.T) {
	type fields struct {
		Name       string
		Addons     map[string]*ClusterAddon
		Properties map[string]string
	}
	type args struct {
		config *ProjectConfig
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantAddons map[string]*ClusterAddon
	}{
		{
			name: "one addon with bool property false",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				config: &ProjectConfig{
					Addons: map[string]Addon{
						"addon1": {
							Group:          "group1",
							DefaultEnabled: true,
						},
					},
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     false,
									Type:        template.PropertyTypeBool,
									Description: "property1",
								},
							},
						},
					},
				},
			},
			wantAddons: map[string]*ClusterAddon{
				"addon1": {
					Enabled:    true,
					Properties: map[string]any{"property1": false},
				},
			},
		},
		{
			name: "one addon with bool property nil",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				config: &ProjectConfig{
					Addons: map[string]Addon{
						"addon1": {
							Group:          "group1",
							DefaultEnabled: true,
						},
					},
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
			},
			wantAddons: map[string]*ClusterAddon{
				"addon1": {
					Enabled:    true,
					Properties: map[string]any{"property1": nil},
				},
			},
		},
		{
			name: "one addon with int property 10",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				config: &ProjectConfig{
					Addons: map[string]Addon{
						"addon1": {
							Group:          "group1",
							DefaultEnabled: true,
						},
					},
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     10,
									Type:        template.PropertyTypeInt,
									Description: "property1",
								},
							},
						},
					},
				},
			},
			wantAddons: map[string]*ClusterAddon{
				"addon1": {
					Enabled:    true,
					Properties: map[string]any{"property1": 10},
				},
			},
		},
		{
			name: "one addon with int property nil",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				config: &ProjectConfig{
					Addons: map[string]Addon{
						"addon1": {
							Group:          "group1",
							DefaultEnabled: true,
						},
					},
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
			},
			wantAddons: map[string]*ClusterAddon{
				"addon1": {
					Enabled:    true,
					Properties: map[string]any{"property1": nil},
				},
			},
		},
		{
			name: "one addon with string property 'Hello World'",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				config: &ProjectConfig{
					Addons: map[string]Addon{
						"addon1": {
							Group:          "group1",
							DefaultEnabled: true,
						},
					},
					ParsedAddons: map[string]template.TemplateManifest{
						"addon1": {
							Properties: map[string]template.Property{
								"property1": {
									Required:    true,
									Default:     "Hello World",
									Type:        template.PropertyTypeString,
									Description: "property1",
								},
							},
						},
					},
				},
			},
			wantAddons: map[string]*ClusterAddon{
				"addon1": {
					Enabled:    true,
					Properties: map[string]any{"property1": "Hello World"},
				},
			},
		},
		{
			name: "one addon with string property 'nil'",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				config: &ProjectConfig{
					Addons: map[string]Addon{
						"addon1": {
							Group:          "group1",
							DefaultEnabled: true,
						},
					},
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
			},
			wantAddons: map[string]*ClusterAddon{
				"addon1": {
					Enabled:    true,
					Properties: map[string]any{"property1": nil},
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
			c.SetDefaultAddons(tt.args.config)

			diff := cmp.Diff(tt.wantAddons, c.Addons)
			if diff != "" {
				t.Errorf("Cluster.SetDefaultAddons() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
