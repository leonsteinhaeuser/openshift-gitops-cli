package utils

import (
	"reflect"
	"slices"
	"testing"
)

func TestMapKeysToList(t *testing.T) {
	type args struct {
		m map[string]any
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "with populated map",
			args: args{
				m: map[string]any{
					"key1": nil,
					"key2": nil,
				},
			},
			want: []string{"key1", "key2"},
		},
		{
			name: "with empty map",
			args: args{
				m: map[string]any{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapKeysToList(tt.args.m)
			// we need to sort the slices because the order of the keys in the map is not guaranteed
			slices.Sort(got)
			slices.Sort(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapKeysToList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	type args struct {
		maps []map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "with populated maps",
			args: args{
				maps: []map[string]string{
					{
						"key1": "value1",
						"key2": "value2",
					},
					{
						"key3": "value3",
					},
				},
			},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "with overwrite maps",
			args: args{
				maps: []map[string]string{
					{
						"key1": "value1",
						"key2": "value2",
					},
					{
						"key1": "value3",
					},
				},
			},
			want: map[string]string{
				"key1": "value3",
				"key2": "value2",
			},
		},
		{
			name: "with empty maps",
			args: args{
				maps: []map[string]string{},
			},
			want: map[string]string{},
		},
		{
			name: "with nil maps",
			args: args{
				maps: []map[string]string{nil},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeMaps(tt.args.maps...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduceMap(t *testing.T) {
	type args struct {
		origin map[string]string
		maps   []map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "with populated maps",
			args: args{
				origin: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				maps: []map[string]string{
					{
						"key3": "value3",
					},
				},
			},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "with populated maps",
			args: args{
				origin: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				maps: []map[string]string{
					{
						"key2": "value3",
					},
				},
			},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "with populated maps",
			args: args{
				origin: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				maps: []map[string]string{
					{
						"key1": "value1",
					},
					{
						"key2": "value2",
					},
					{
						"key1": "value1",
					},
				},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReduceMap(tt.args.origin, tt.args.maps...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReduceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
