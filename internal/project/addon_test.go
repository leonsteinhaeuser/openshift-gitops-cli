package project

import (
	"testing"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
)

func TestClusterAddons_AllRequiredPropertiesSet(t *testing.T) {
	type fields struct {
		Addons ClusterAddons
	}
	type args struct {
		config        *ProjectConfig
		skipOnFailure bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no addons defined",
			fields: fields{
				Addons: ClusterAddons{},
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
				skipOnFailure: false,
			},
			wantErr: false,
		},
		{
			name: "one property set",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": true,
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: false,
		},
		{
			name: "one property set string, expect bool",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": "invalid",
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect bool",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": 10,
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: true,
		},
		{
			name: "one property set bool, expect string",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": true,
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect string",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": 10,
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: true,
		},
		{
			name: "one property set string, expect string",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": "value",
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: false,
		},
		{
			name: "one property set int, expect int",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": 10,
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: false,
		},
		{
			name: "one property set string, expect int",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": "value",
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect int",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": 10,
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: false,
		},
		{
			name: "one property set int, expect int, invalid value",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": "invalid",
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: true,
		},
		{
			name: "one property set int, expect int, invalid value",
			fields: fields{
				Addons: ClusterAddons{
					"addon1": {
						Enabled: true,
						Properties: map[string]any{
							"property1": "invalid",
						},
					},
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
				skipOnFailure: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := tt.fields.Addons
			if err := ca.AllRequiredPropertiesSet(tt.args.config, tt.args.skipOnFailure); (err != nil) != tt.wantErr {
				t.Errorf("ClusterAddons.AllRequiredPropertiesSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

func TestClusterAddons_IsEnabled(t *testing.T) {
	type args struct {
		addon string
	}
	tests := []struct {
		name string
		ca   ClusterAddons
		args args
		want bool
	}{
		{
			name: "addon not enabled",
			ca: ClusterAddons{
				"addon1": {
					Enabled: false,
				},
			},
			args: args{
				addon: "addon1",
			},
			want: false,
		},
		{
			name: "addon enabled",
			ca: ClusterAddons{
				"addon1": {
					Enabled: true,
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
			if got := tt.ca.IsEnabled(tt.args.addon); got != tt.want {
				t.Errorf("ClusterAddons.IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterAddon_IsEnabled(t *testing.T) {
	type fields struct {
		Enabled    bool
		Properties map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "addon not enabled",
			fields: fields{
				Enabled: false,
			},
			want: false,
		},
		{
			name: "addon enabled",
			fields: fields{
				Enabled: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := ClusterAddon{
				Enabled:    tt.fields.Enabled,
				Properties: tt.fields.Properties,
			}
			if got := ca.IsEnabled(); got != tt.want {
				t.Errorf("ClusterAddon.IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterAddon_SetProperty(t *testing.T) {
	type fields struct {
		Enabled    bool
		Properties map[string]any
	}
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set property",
			fields: fields{
				Properties: map[string]any{},
			},
			args: args{
				key:   "property1",
				value: "value",
			},
		},
		{
			name: "set property, override",
			fields: fields{
				Properties: map[string]any{
					"property1": "value",
				},
			},
			args: args{
				key:   "property1",
				value: "value2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := &ClusterAddon{
				Enabled:    tt.fields.Enabled,
				Properties: tt.fields.Properties,
			}
			ca.SetProperty(tt.args.key, tt.args.value)
		})
	}
}
