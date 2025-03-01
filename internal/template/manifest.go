package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

/*
// Directory structure

- manifest.yaml     # contains the lookup parameters for the template
- values.yaml       # contains actual template related values
- files/abc.txt     # contains actual template related values
- files/def.txt     # contains actual template related values
- files/ghi.yaml    # contains actual template related value

// In config

name: "my-template"
properties:
  globalServerSideApply:
    required: false
	default: false
	description: "Whether to use server-side apply for all resources"
files:
  - values.yaml
  - files/abc.txt
  - files/def.txt
  - files/ghi.yaml
*/

type TemplateManifest struct {
	BasePath string `json:"-"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	// Properties is a map of properties that can be set when rendering the template
	// The key is the name of the property, the value is the property definition
	// In the manifest file, the properties can be defined to be required or have a default value
	// The description key will be used to display a description of the property
	Properties map[string]Property `json:"properties"`
	// Annotations is a map of annotations that can be set when rendering the template
	// The key is the name of the annotation, the value is the annotation definition
	Annotations map[string]string `json:"annotations"`
	// Files is a list of relative paths to files that are part of the template
	Files []string `json:"files"`
}

type PropertyType string

const (
	PropertyTypeString PropertyType = "string"
	PropertyTypeBool   PropertyType = "bool"
	PropertyTypeInt    PropertyType = "int"
)

// checkType validates the given value against the property type
// If the value is valid, it will be returned, otherwise an error is returned
func (p PropertyType) checkType(value any) (any, error) {
	switch v := value.(type) {
	case string:
		if p != PropertyTypeString {
			return nil, fmt.Errorf("expected type %s, got string", p)
		}
		return v, nil
	case bool:
		if p != PropertyTypeBool {
			return nil, fmt.Errorf("expected type %s, got bool", p)
		}
		return v, nil
	case int:
		if p != PropertyTypeInt {
			return nil, fmt.Errorf("expected type %s, got int", p)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}

type Property struct {
	Required    bool
	Default     any
	Type        PropertyType
	Description string
}

// Check validates the given value against the property definition
func (p Property) ParseValue(value any) (any, error) {
	v := p.Default
	if value != nil {
		v = value
	}
	if p.Required && v == nil {
		return nil, fmt.Errorf("property is required")
	}
	if v == nil {
		return nil, nil
	}
	return p.Type.checkType(v)
}

// LoadManifest reads the manifest file at the given path and returns the parsed values
// If the path is a directory, it will try to find the manifest file in it
func LoadManifest(path string) (*TemplateManifest, error) {
	// if user provides a directory, try to find the manifest file in it
	if !strings.HasSuffix(path, "manifest.yaml") && !strings.HasSuffix(path, "manifest.yml") {
		path = filepath.Join(path, "manifest.yaml")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = strings.TrimSuffix(path, "manifest.yaml") + "manifest.yml"
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("no manifest file found in %s", path)
		}
	}
	// read manifest file and parse it
	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t := &TemplateManifest{}
	err = yaml.Unmarshal(bts, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// LoadTemplateManifest walks the given basePath and loads all template manifests files it finds
// The returned slice contains the path to the manifest file and the parsed values
func LoadTemplateManifest(basePath string) ([]Template, error) {
	templates := []Template{}
	err := filepath.WalkDir(basePath, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// skip directories
			return nil
		}

		if d.Name() != "manifest.yaml" && d.Name() != "manifest.yml" {
			// not a manifest file
			return nil
		}

		// now read the manifest file
		tm, err := LoadManifest(fpath)
		if err != nil {
			return err
		}

		templates = append(templates, Template{
			Path:             filepath.Dir(fpath),
			TemplateManifest: *tm,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}
