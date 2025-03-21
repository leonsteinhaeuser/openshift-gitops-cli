package project

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStage_IsAddonEnabled(t *testing.T) {
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
			s := Stage{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			if got := s.IsAddonEnabled(tt.args.addon); got != tt.want {
				t.Errorf("Stage.IsAddonEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStage_EnableAddon(t *testing.T) {
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
		want   Stage
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				addon: "addon1",
			},
			want: Stage{
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
			want: Stage{
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
			want: Stage{
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
			want: Stage{
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
			want: Stage{
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
			s := &Stage{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			s.EnableAddon(tt.args.addon)

			diff := cmp.Diff(tt.want, *s)
			if diff != "" {
				t.Errorf("Stage.EnableAddon() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestStage_DisableAddon(t *testing.T) {
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
		want   Stage
	}{
		{
			name: "no addons",
			fields: fields{
				Addons: map[string]*ClusterAddon{},
			},
			args: args{
				addon: "addon1",
			},
			want: Stage{
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
			want: Stage{
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
			want: Stage{
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
			s := &Stage{
				Name:       tt.fields.Name,
				Addons:     tt.fields.Addons,
				Properties: tt.fields.Properties,
			}
			s.DisableAddon(tt.args.addon)

			diff := cmp.Diff(tt.want, *s)
			if diff != "" {
				t.Errorf("Stage.DisableAddon() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestStage_GetAddons(t *testing.T) {
	type fields struct {
		Name       string
		Properties map[string]string
		Actions    Actions
		Clusters   map[string]*Cluster
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
			s := &Stage{
				Name:       tt.fields.Name,
				Properties: tt.fields.Properties,
				Actions:    tt.fields.Actions,
				Clusters:   tt.fields.Clusters,
				Addons:     tt.fields.Addons,
			}
			if got := s.GetAddons(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stage.GetAddons() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStage_GetAddon(t *testing.T) {
	type fields struct {
		Name       string
		Properties map[string]string
		Actions    Actions
		Clusters   map[string]*Cluster
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
			s := &Stage{
				Name:       tt.fields.Name,
				Properties: tt.fields.Properties,
				Actions:    tt.fields.Actions,
				Clusters:   tt.fields.Clusters,
				Addons:     tt.fields.Addons,
			}
			if got := s.GetAddon(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stage.GetAddon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStage_GetCluster(t *testing.T) {
	type fields struct {
		Name       string
		Properties map[string]string
		Actions    Actions
		Clusters   map[string]*Cluster
		Addons     map[string]*ClusterAddon
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Cluster
	}{
		{
			name: "no clusters",
			fields: fields{
				Clusters: map[string]*Cluster{},
			},
			args: args{
				name: "cluster1",
			},
			want: nil,
		},
		{
			name: "one cluster",
			fields: fields{
				Clusters: map[string]*Cluster{
					"cluster1": {
						Name: "cluster1",
					},
				},
			},
			args: args{
				name: "cluster1",
			},
			want: &Cluster{
				Name: "cluster1",
			},
		},
		{
			name: "two clusters",
			fields: fields{
				Clusters: map[string]*Cluster{
					"cluster1": {
						Name: "cluster1",
					},
					"cluster2": {
						Name: "cluster2",
					},
				},
			},
			args: args{
				name: "cluster1",
			},
			want: &Cluster{
				Name: "cluster1",
			},
		},
		{
			name: "two clusters no found",
			fields: fields{
				Clusters: map[string]*Cluster{
					"cluster1": {
						Name: "cluster1",
					},
					"cluster2": {
						Name: "cluster2",
					},
				},
			},
			args: args{
				name: "cluster5",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stage{
				Name:       tt.fields.Name,
				Properties: tt.fields.Properties,
				Actions:    tt.fields.Actions,
				Clusters:   tt.fields.Clusters,
				Addons:     tt.fields.Addons,
			}
			if got := s.GetCluster(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stage.GetCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}
