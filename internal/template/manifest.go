package template

import (
	"io/fs"
	"os"
	"path/filepath"

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
	Name string
	// Properties is a map of properties that can be set when rendering the template
	// The key is the name of the property, the value is the property definition
	// In the manifest file, the properties can be defined to be required or have a default value
	// The description key will be used to display a description of the property
	Properties map[string]Property
	// Files is a list of relative paths to files that are part of the template
	Files []string
}

type Property struct {
	Required    bool
	Default     string
	Description string
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
		bts, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}
		t := TemplateManifest{}
		err = yaml.Unmarshal(bts, &t)
		if err != nil {
			return err
		}

		templates = append(templates, Template{
			Path:             filepath.Dir(fpath),
			TemplateManifest: t,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}
