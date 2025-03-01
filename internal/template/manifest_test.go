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

func TestProperty_ParseValue(t *testing.T) {
	type fields struct {
		Required    bool
		Default     any
		Type        PropertyType
		Description string
	}
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// string
		{
			name: "string type with string value",
			fields: fields{
				Type: "string",
			},
			args: args{
				value: "test",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "string type with default string value",
			fields: fields{
				Type:    "string",
				Default: "default",
			},
			args: args{
				value: nil,
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "string type with default string and string value",
			fields: fields{
				Type:    "string",
				Default: "default",
			},
			args: args{
				value: "test",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "string type with required value",
			fields: fields{
				Type:     "string",
				Required: true,
			},
			args: args{
				value: "test",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "string type with required value",
			fields: fields{
				Type:     "string",
				Required: true,
			},
			args: args{
				value: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "string type with int value",
			fields: fields{
				Type: "string",
			},
			args: args{
				value: 42,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "string type with bool value",
			fields: fields{
				Type: "string",
			},
			args: args{
				value: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "string type with nil value",
			fields: fields{
				Type: "string",
			},
			args: args{
				value: nil,
			},
			want:    nil,
			wantErr: false,
		},
		// int
		{
			name: "int type with int value",
			fields: fields{
				Type: "int",
			},
			args: args{
				value: 10,
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "int type with default int value",
			fields: fields{
				Type:    "int",
				Default: 42,
			},
			args: args{
				value: nil,
			},
			want:    42,
			wantErr: false,
		},
		{
			name: "int type with default int and int value",
			fields: fields{
				Type:    "int",
				Default: 42,
			},
			args: args{
				value: 12,
			},
			want:    12,
			wantErr: false,
		},
		{
			name: "int type with required value",
			fields: fields{
				Type:     "int",
				Required: true,
			},
			args: args{
				value: 2,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "int type with required value",
			fields: fields{
				Type:     "int",
				Required: true,
			},
			args: args{
				value: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "int type with string value",
			fields: fields{
				Type: "int",
			},
			args: args{
				value: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "int type with bool value",
			fields: fields{
				Type: "int",
			},
			args: args{
				value: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "int type with nil value",
			fields: fields{
				Type: "int",
			},
			args: args{
				value: nil,
			},
			want:    nil,
			wantErr: false,
		},
		// bool
		{
			name: "bool type with bool value",
			fields: fields{
				Type: "bool",
			},
			args: args{
				value: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "bool type with default bool value",
			fields: fields{
				Type:    "bool",
				Default: true,
			},
			args: args{
				value: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "bool type with default bool and bool value",
			fields: fields{
				Type:    "bool",
				Default: true,
			},
			args: args{
				value: false,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "bool type with required value",
			fields: fields{
				Type:     "bool",
				Required: true,
			},
			args: args{
				value: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "bool type with required value",
			fields: fields{
				Type:     "bool",
				Required: true,
			},
			args: args{
				value: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bool type with string value",
			fields: fields{
				Type: "bool",
			},
			args: args{
				value: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bool type with int value",
			fields: fields{
				Type: "bool",
			},
			args: args{
				value: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bool type with nil value",
			fields: fields{
				Type: "int",
			},
			args: args{
				value: nil,
			},
			want:    nil,
			wantErr: false,
		},
		// unknown type
		{
			name: "unknown type",
			fields: fields{
				Type: "unknown",
			},
			args: args{
				value: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "complex value",
			fields: fields{
				Type: "string",
			},
			args: args{
				value: []string{"test"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Property{
				Required:    tt.fields.Required,
				Default:     tt.fields.Default,
				Type:        tt.fields.Type,
				Description: tt.fields.Description,
			}
			got, err := p.ParseValue(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Property.ParseValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Property.ParseValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
