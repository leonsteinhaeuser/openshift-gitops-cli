package utils

import (
	"reflect"
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
			if got := MapKeysToList(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapKeysToList() = %v, want %v", got, tt.want)
			}
		})
	}
}
