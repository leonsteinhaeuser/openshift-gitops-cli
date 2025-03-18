package project

import (
	"reflect"
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

func TestEnvironment_GetAddons(t *testing.T) {
	type fields struct {
		Name       string
		Properties map[string]string
		Actions    Actions
		Stages     map[string]*Stage
		Addons     map[string]*ClusterAddon
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
			e := &Environment{
				Name:       tt.fields.Name,
				Properties: tt.fields.Properties,
				Actions:    tt.fields.Actions,
				Stages:     tt.fields.Stages,
				Addons:     tt.fields.Addons,
			}
			if got := e.GetAddons(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Environment.GetAddons() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_GetAddon(t *testing.T) {
	type fields struct {
		Name       string
		Properties map[string]string
		Actions    Actions
		Stages     map[string]*Stage
		Addons     map[string]*ClusterAddon
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
			e := &Environment{
				Name:       tt.fields.Name,
				Properties: tt.fields.Properties,
				Actions:    tt.fields.Actions,
				Stages:     tt.fields.Stages,
				Addons:     tt.fields.Addons,
			}
			if got := e.GetAddon(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Environment.GetAddon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_HasStage(t *testing.T) {
	type fields struct {
		Name       string
		Properties map[string]string
		Actions    Actions
		Stages     map[string]*Stage
		Addons     map[string]*ClusterAddon
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "no stages",
			fields: fields{
				Stages: map[string]*Stage{},
			},
			args: args{
				name: "stage1",
			},
			want: false,
		},
		{
			name: "one stage",
			fields: fields{
				Stages: map[string]*Stage{
					"stage1": {},
				},
			},
			args: args{
				name: "stage1",
			},
			want: true,
		},
		{
			name: "two stages and found",
			fields: fields{
				Stages: map[string]*Stage{
					"stage1": {},
					"stage2": {},
				},
			},
			args: args{
				name: "stage1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Environment{
				Name:       tt.fields.Name,
				Properties: tt.fields.Properties,
				Actions:    tt.fields.Actions,
				Stages:     tt.fields.Stages,
				Addons:     tt.fields.Addons,
			}
			if got := e.HasStage(tt.args.name); got != tt.want {
				t.Errorf("Environment.HasStage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_GetStage(t *testing.T) {
	type fields struct {
		Name       string
		Properties map[string]string
		Actions    Actions
		Stages     map[string]*Stage
		Addons     map[string]*ClusterAddon
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Stage
	}{
		{
			name: "no stages",
			fields: fields{
				Stages: map[string]*Stage{},
			},
			args: args{
				name: "stage1",
			},
			want: nil,
		},
		{
			name: "one stage",
			fields: fields{
				Stages: map[string]*Stage{
					"stage1": {},
				},
			},
			args: args{
				name: "stage1",
			},
			want: &Stage{},
		},
		{
			name: "two stages and found",
			fields: fields{
				Stages: map[string]*Stage{
					"stage1": {},
					"stage2": {},
				},
			},
			args: args{
				name: "stage1",
			},
			want: &Stage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Environment{
				Name:       tt.fields.Name,
				Properties: tt.fields.Properties,
				Actions:    tt.fields.Actions,
				Stages:     tt.fields.Stages,
				Addons:     tt.fields.Addons,
			}
			if got := e.GetStage(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Environment.GetStage() = %v, want %v", got, tt.want)
			}
		})
	}
}
