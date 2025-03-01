package template

import "testing"

func Test_toYAML(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with multiple keys",
			args: args{
				v: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: "key1: value1\nkey2: value2",
		},
		{
			name: "with one key",
			args: args{
				v: map[string]string{
					"key1": "value1",
				},
			},
			want: "key1: value1",
		},
		{
			name: "with empty map",
			args: args{
				v: map[string]string{},
			},
			want: "{}", // if empty, the function returns an empty object {}
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toYAML(tt.args.v); got != tt.want {
				t.Errorf("toYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}
