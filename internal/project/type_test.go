package project

import (
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
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ProjectConfig.AddonGroups() mismatch (-got +want):\n%s", diff)
				return
			}
		})
	}
}
