package template

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"strings"

	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/yaml"
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
				value: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "int",
			p:    PropertyTypeInt,
			args: args{
				value: 42,
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
		{
			name: "property is bool but expect string",
			p:    PropertyTypeString,
			args: args{
				value: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unsupported data",
			p:    PropertyTypeString,
			args: args{
				value: []string{"test"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil value",
			p:    PropertyTypeString,
			args: args{
				value: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.checkType(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("PropertyType.checkType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("PropertyType.checkType() mismatch (-want +got):\n%s", diff)
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

func TestLoadManifest(t *testing.T) {
	type args struct {
		path string
	}
	type actions struct {
		PreTest  func(path string)
		PostTest func(path string)
	}
	tests := []struct {
		name    string
		action  actions
		args    args
		want    *TemplateManifest
		wantErr bool
	}{
		{
			name: "without filename in path",
			action: actions{
				PreTest: func(path string) {
					file, err := os.OpenFile(filepath.Join(path, "manifest.yaml"), os.O_CREATE|os.O_WRONLY, 0664)
					if err != nil {
						t.Fatal(err)
						return
					}
					defer file.Close()

					bts, err := yaml.Marshal(&TemplateManifest{
						Name:  "test",
						Group: "test",
						Properties: map[string]Property{
							"test": {
								Required:    true,
								Default:     "default",
								Type:        PropertyTypeString,
								Description: "test",
							},
						},
						Annotations: map[string]string{
							"test": "test",
						},
						Files: []string{"test"},
					})
					if err != nil {
						t.Fatal(err)
						return
					}

					_, err = file.Write(bts)
					if err != nil {
						t.Fatal(err)
						return
					}
				},
				PostTest: func(path string) {
					err := os.Remove(filepath.Join(path, "manifest.yaml"))
					if err != nil {
						t.Fatal(err)
						return
					}
				},
			},
			args: args{
				path: os.TempDir(),
			},
			want: &TemplateManifest{
				BasePath: "test",
				Name:     "test",
				Group:    "test",
				Properties: map[string]Property{
					"test": {
						Required:    true,
						Default:     "default",
						Type:        PropertyTypeString,
						Description: "test",
					},
				},
				Annotations: map[string]string{
					"test": "test",
				},
				Files: []string{"test"},
			},
			wantErr: false,
		},
		{
			name: "with filename in path",
			action: actions{
				PreTest: func(path string) {
					file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0664)
					if err != nil {
						t.Fatal(err)
						return
					}
					defer file.Close()

					bts, err := yaml.Marshal(&TemplateManifest{
						Name:  "test",
						Group: "test",
						Properties: map[string]Property{
							"test": {
								Required:    true,
								Default:     "default",
								Type:        PropertyTypeString,
								Description: "test",
							},
						},
						Annotations: map[string]string{
							"test": "test",
						},
						Files: []string{"test"},
					})
					if err != nil {
						t.Fatal(err)
						return
					}

					_, err = file.Write(bts)
					if err != nil {
						t.Fatal(err)
						return
					}
				},
				PostTest: func(path string) {
					err := os.Remove(path)
					if err != nil {
						t.Fatal(err)
						return
					}
				},
			},
			args: args{
				path: filepath.Join(os.TempDir(), "manifest.yaml"),
			},
			want: &TemplateManifest{
				BasePath: "test",
				Name:     "test",
				Group:    "test",
				Properties: map[string]Property{
					"test": {
						Required:    true,
						Default:     "default",
						Type:        PropertyTypeString,
						Description: "test",
					},
				},
				Annotations: map[string]string{
					"test": "test",
				},
				Files: []string{"test"},
			},
			wantErr: false,
		},
		{
			name: "file not exists",
			action: actions{
				PreTest:  func(path string) {},
				PostTest: func(path string) {},
			},
			args: args{
				path: os.TempDir(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.action.PreTest(tt.args.path)
			defer tt.action.PostTest(tt.args.path)

			got, err := LoadManifest(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("LoadManifest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func mockTemplateFile(path string, templateName string) error {
	// create root directory
	path = filepath.Join(path, templateName)
	err := os.MkdirAll(path, 0775)
	if err != nil {
		return err
	}

	// create a file in a subdirectory
	configDir := filepath.Join(path, "config")
	err = os.MkdirAll(configDir, 0775)
	if err != nil {
		return err
	}

	// create a file in a subdirectory
	patchesDir := filepath.Join(path, "patches")
	err = os.MkdirAll(patchesDir, 0775)
	if err != nil {
		return err
	}

	// create the addon entrypoint file
	file, err := os.OpenFile(filepath.Join(path, "manifest.yaml"), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	bts, err := yaml.Marshal(&TemplateManifest{
		Name:  templateName,
		Group: templateName,
		Properties: map[string]Property{
			"test": {
				Required:    true,
				Default:     "default",
				Type:        PropertyTypeString,
				Description: "test",
			},
		},
		Annotations: map[string]string{
			"test": "test",
		},
		Files: []string{"./"},
	})
	if err != nil {
		return err
	}

	_, err = file.Write(bts)
	if err != nil {
		return err
	}

	//
	// create some example files
	//
	file, err = os.OpenFile(filepath.Join(path, "test.yaml"), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte("test: {{ .Properties.test }}"))
	if err != nil {
		return err
	}

	file, err = os.OpenFile(filepath.Join(configDir, "test.yaml"), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte("test: {{ .Properties.test }}"))
	if err != nil {
		return err
	}

	file, err = os.OpenFile(filepath.Join(patchesDir, "test.yaml"), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte("test: {{ .Properties.test }}"))
	if err != nil {
		return err
	}

	return nil
}

func TestLoadTemplateManifest(t *testing.T) {
	type args struct {
		path string
	}
	type actions struct {
		PreTest func(path string)
	}
	tests := []struct {
		name    string
		action  actions
		args    args
		want    []Template
		wantErr bool
	}{
		{
			name: "without filename in path",
			action: actions{
				PreTest: func(path string) {
					err := mockTemplateFile(path, "one")
					if err != nil {
						t.Fatal(err)
						return
					}

					err = mockTemplateFile(path, "two")
					if err != nil {
						t.Fatal(err)
						return
					}
				},
			},
			args: args{
				path: os.TempDir(),
			},
			want: []Template{
				{
					Path: "one",
					TemplateManifest: TemplateManifest{
						BasePath: "one",
						Name:     "one",
						Group:    "one",
						Properties: map[string]Property{
							"test": {
								Required:    true,
								Default:     "default",
								Type:        PropertyTypeString,
								Description: "test",
							},
						},
						Annotations: map[string]string{
							"test": "test",
						},
						Files: []string{"./"},
					},
				},
				{
					Path: "two",
					TemplateManifest: TemplateManifest{
						BasePath: "two",
						Name:     "two",
						Group:    "two",
						Properties: map[string]Property{
							"test": {
								Required:    true,
								Default:     "default",
								Type:        PropertyTypeString,
								Description: "test",
							},
						},
						Annotations: map[string]string{
							"test": "test",
						},
						Files: []string{"./"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			tt.action.PreTest(dir)

			got, err := LoadTemplateManifest(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTemplateManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// in order to compare the paths, we need to sanitize the paths, to make it easier for the comparison
			for idx, tmpl := range got {
				tmpl.Path = strings.TrimPrefix(strings.TrimPrefix(tmpl.Path, dir), string(os.PathSeparator))
				got[idx] = tmpl
			}

			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("LoadTemplateManifest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
