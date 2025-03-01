package template

import (
	"reflect"
	"testing"
)

func TestPropertyType_checkType(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		p       PropertyType
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "string",
			p:    PropertyTypeString,
			args: args{
				value: "test",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "bool",
			p:    PropertyTypeBool,
			args: args{
				value: "true",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "int",
			p:    PropertyTypeInt,
			args: args{
				value: "42",
			},
			want:    int64(42),
			wantErr: false,
		},
		{
			name: "invalid string",
			p:    PropertyTypeString,
			args: args{
				value: 42,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid bool",
			p:    PropertyTypeBool,
			args: args{
				value: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid int",
			p:    PropertyTypeInt,
			args: args{
				value: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown type",
			p:    "unknown",
			args: args{
				value: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.checkType(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("PropertyType.checkType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PropertyType.checkType() = %v, want %v", got, tt.want)
			}
		})
	}
}
