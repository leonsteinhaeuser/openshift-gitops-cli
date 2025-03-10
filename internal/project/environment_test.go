package project

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvironment_IsAddonEnabled(t *testing.T) {
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
			e := Environment{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			if got := e.IsAddonEnabled(tt.args.addon); got != tt.want {
				t.Errorf("Environment.IsAddonEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_EnableAddon(t *testing.T) {
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
		want   Environment
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				addon: "addon1",
			},
			want: Environment{
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
			want: Environment{
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
			want: Environment{
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
			want: Environment{
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
			want: Environment{
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
			e := &Environment{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			e.EnableAddon(tt.args.addon)

			diff := cmp.Diff(tt.want, *e)
			if diff != "" {
				t.Errorf("Environment.EnableAddon() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestEnvironment_DisableAddon(t *testing.T) {
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
		want   Environment
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				addon: "addon1",
			},
			want: Environment{
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
			want: Environment{
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
			want: Environment{
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
			e := &Environment{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			e.DisableAddon(tt.args.addon)

			diff := cmp.Diff(tt.want, *e)
			if diff != "" {
				t.Errorf("Environment.DisableAddon() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
