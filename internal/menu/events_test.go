package menu

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_newPreCreateEvent(t *testing.T) {
	type args struct {
		origin      EventOrigin
		environment string
		stage       string
		cluster     string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "with environment",
			args: args{
				origin:      EventOriginEnvironment,
				environment: "test",
			},
			want: Event{
				Type:        EventTypeCreate,
				Origin:      EventOriginEnvironment,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "",
				Cluster:     "",
			},
		},
		{
			name: "with environment and stage",
			args: args{
				origin:      EventOriginStage,
				environment: "test",
				stage:       "test",
			},
			want: Event{
				Type:        EventTypeCreate,
				Origin:      EventOriginStage,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "test",
				Cluster:     "",
			},
		},
		{
			name: "with environment, stage and cluster",
			args: args{
				origin:      EventOriginCluster,
				environment: "test",
				stage:       "test",
				cluster:     "test",
			},
			want: Event{
				Type:        EventTypeCreate,
				Origin:      EventOriginCluster,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "test",
				Cluster:     "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newPreCreateEvent(tt.args.origin, tt.args.environment, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("newPreCreateEvent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_newPostCreateEvent(t *testing.T) {
	type args struct {
		origin      EventOrigin
		environment string
		stage       string
		cluster     string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "with environment",
			args: args{
				origin:      EventOriginEnvironment,
				environment: "test",
			},
			want: Event{
				Type:        EventTypeCreate,
				Origin:      EventOriginEnvironment,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "",
				Cluster:     "",
			},
		},
		{
			name: "with environment and stage",
			args: args{
				origin:      EventOriginStage,
				environment: "test",
				stage:       "test",
			},
			want: Event{
				Type:        EventTypeCreate,
				Origin:      EventOriginStage,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "test",
				Cluster:     "",
			},
		},
		{
			name: "with environment, stage and cluster",
			args: args{
				origin:      EventOriginCluster,
				environment: "test",
				stage:       "test",
				cluster:     "test",
			},
			want: Event{
				Type:        EventTypeCreate,
				Origin:      EventOriginCluster,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "test",
				Cluster:     "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newPostCreateEvent(tt.args.origin, tt.args.environment, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("newPostCreateEvent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_newPreUpdateEvent(t *testing.T) {
	type args struct {
		origin      EventOrigin
		environment string
		stage       string
		cluster     string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "with environment",
			args: args{
				origin:      EventOriginEnvironment,
				environment: "test",
			},
			want: Event{
				Type:        EventTypeUpdate,
				Origin:      EventOriginEnvironment,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "",
				Cluster:     "",
			},
		},
		{
			name: "with environment and stage",
			args: args{
				origin:      EventOriginStage,
				environment: "test",
				stage:       "test",
			},
			want: Event{
				Type:        EventTypeUpdate,
				Origin:      EventOriginStage,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "test",
				Cluster:     "",
			},
		},
		{
			name: "with environment, stage and cluster",
			args: args{
				origin:      EventOriginCluster,
				environment: "test",
				stage:       "test",
				cluster:     "test",
			},
			want: Event{
				Type:        EventTypeUpdate,
				Origin:      EventOriginCluster,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "test",
				Cluster:     "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newPreUpdateEvent(tt.args.origin, tt.args.environment, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("newPreUpdateEvent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_newPostUpdateEvent(t *testing.T) {
	type args struct {
		origin      EventOrigin
		environment string
		stage       string
		cluster     string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "with environment",
			args: args{
				origin:      EventOriginEnvironment,
				environment: "test",
			},
			want: Event{
				Type:        EventTypeUpdate,
				Origin:      EventOriginEnvironment,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "",
				Cluster:     "",
			},
		},
		{
			name: "with environment and stage",
			args: args{
				origin:      EventOriginStage,
				environment: "test",
				stage:       "test",
			},
			want: Event{
				Type:        EventTypeUpdate,
				Origin:      EventOriginStage,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "test",
				Cluster:     "",
			},
		},
		{
			name: "with environment, stage and cluster",
			args: args{
				origin:      EventOriginCluster,
				environment: "test",
				stage:       "test",
				cluster:     "test",
			},
			want: Event{
				Type:        EventTypeUpdate,
				Origin:      EventOriginCluster,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "test",
				Cluster:     "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newPostUpdateEvent(tt.args.origin, tt.args.environment, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("newPostUpdateEvent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_newPreDeleteEvent(t *testing.T) {
	type args struct {
		origin      EventOrigin
		environment string
		stage       string
		cluster     string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "with environment",
			args: args{
				origin:      EventOriginEnvironment,
				environment: "test",
			},
			want: Event{
				Type:        EventTypeDelete,
				Origin:      EventOriginEnvironment,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "",
				Cluster:     "",
			},
		},
		{
			name: "with environment and stage",
			args: args{
				origin:      EventOriginStage,
				environment: "test",
				stage:       "test",
			},
			want: Event{
				Type:        EventTypeDelete,
				Origin:      EventOriginStage,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "test",
				Cluster:     "",
			},
		},
		{
			name: "with environment, stage and cluster",
			args: args{
				origin:      EventOriginCluster,
				environment: "test",
				stage:       "test",
				cluster:     "test",
			},
			want: Event{
				Type:        EventTypeDelete,
				Origin:      EventOriginCluster,
				Runtime:     EventRuntimePre,
				Environment: "test",
				Stage:       "test",
				Cluster:     "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newPreDeleteEvent(tt.args.origin, tt.args.environment, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("newPreDeleteEvent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_newPostDeleteEvent(t *testing.T) {
	type args struct {
		origin      EventOrigin
		environment string
		stage       string
		cluster     string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "with environment",
			args: args{
				origin:      EventOriginEnvironment,
				environment: "test",
			},
			want: Event{
				Type:        EventTypeDelete,
				Origin:      EventOriginEnvironment,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "",
				Cluster:     "",
			},
		},
		{
			name: "with environment and stage",
			args: args{
				origin:      EventOriginStage,
				environment: "test",
				stage:       "test",
			},
			want: Event{
				Type:        EventTypeDelete,
				Origin:      EventOriginStage,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "test",
				Cluster:     "",
			},
		},
		{
			name: "with environment, stage and cluster",
			args: args{
				origin:      EventOriginCluster,
				environment: "test",
				stage:       "test",
				cluster:     "test",
			},
			want: Event{
				Type:        EventTypeDelete,
				Origin:      EventOriginCluster,
				Runtime:     EventRuntimePost,
				Environment: "test",
				Stage:       "test",
				Cluster:     "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newPostDeleteEvent(tt.args.origin, tt.args.environment, tt.args.stage, tt.args.cluster)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("newPostDeleteEvent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
