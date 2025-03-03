package menu

import (
	"testing"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
)

func Test_allRequiredPropertiesSet(t *testing.T) {
	type args struct {
		addonConfig template.TemplateManifest
		properties  map[string]any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all required properties set",
			args: args{
				addonConfig: template.TemplateManifest{
					Properties: map[string]template.Property{
						"prop1": {
							Required: true,
							Type:     "string",
						},
						"prop2": {
							Required: true,
							Type:     "string",
						},
					},
				},
				properties: map[string]any{
					"prop1": "value1",
					"prop2": "value2",
				},
			},
			want: true,
		},
		{
			name: "not all required properties set",
			args: args{
				addonConfig: template.TemplateManifest{
					Properties: map[string]template.Property{
						"prop1": {
							Required: true,
							Type:     "string",
						},
						"prop2": {
							Required: true,
							Type:     "string",
						},
					},
				},
				properties: map[string]any{
					"prop1": "value1",
				},
			},
			want: false,
		},
		{
			name: "no required properties set",
			args: args{
				addonConfig: template.TemplateManifest{
					Properties: map[string]template.Property{
						"prop1": {
							Required: true,
							Type:     "string",
						},
						"prop2": {
							Required: true,
							Type:     "string",
						},
					},
				},
				properties: map[string]any{},
			},
			want: false,
		},
		{
			name: "no required properties set but not all required",
			args: args{
				addonConfig: template.TemplateManifest{
					Properties: map[string]template.Property{
						"prop1": {
							Required: true,
							Type:     "string",
						},
						"prop2": {
							Required: false,
							Type:     "string",
						},
					},
				},
				properties: map[string]any{},
			},
			want: false,
		},
		{
			name: "no required properties set but not all required",
			args: args{
				addonConfig: template.TemplateManifest{
					Properties: map[string]template.Property{
						"prop1": {
							Required: true,
							Type:     "string",
						},
						"prop2": {
							Required: false,
							Type:     "string",
						},
					},
				},
				properties: map[string]any{
					"prop2": "value2",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allRequiredPropertiesSet(tt.args.addonConfig, tt.args.properties); got != tt.want {
				t.Errorf("allRequiredPropertiesSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
