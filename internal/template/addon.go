package template

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

type AddonTemplateCarrier struct {
	Name  string
	Group string
	// Files represents a map of file names and their respective templates
	Files map[string]*template.Template
}

func LoadTemplatesFromAddonManifest(source TemplateManifest) (*AddonTemplateCarrier, error) {
	template := &AddonTemplateCarrier{
		Name:  source.Name,
		Group: source.Group,
		Files: map[string]*template.Template{},
	}
	err := filepath.WalkDir(source.BasePath, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// skip directories
			return nil
		}

		// ignore manifest files in addon templates
		if strings.HasSuffix(fpath, "manifest.yaml") || strings.HasSuffix(fpath, "manifest.yml") {
			// skip manifest files
			return nil
		}

		bPath := filepath.Base(fpath)
		isFound := slices.IndexFunc(source.Files, func(indexEntry string) bool {
			return indexEntry == bPath
		})
		if isFound == -1 {
			// not a file that is part of the template
			return nil
		}

		tmpl, err := parseFile(fpath)
		if err != nil {
			return err
		}

		template.Files[bPath] = tmpl
		return nil
	})
	if err != nil {
		return nil, err
	}
	return template, nil
}

type AddonTemplateData struct {
	Environment string
	Stage       string
	Cluster     string
	Properties  map[string]any
}

func (a AddonTemplateCarrier) Render(basePath string, properties AddonTemplateData) error {
	if len(a.Files) == 0 {
		// nothing to render
		return nil
	}
	originPath := path.Join(basePath, properties.Environment, properties.Stage, properties.Cluster, a.Group, a.Name)
	err := os.MkdirAll(originPath, 0775)
	if err != nil {
		return err
	}
	for fileName, template := range a.Files {
		// create the file and render the template
		file, err := os.OpenFile(path.Join(originPath, fileName), os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
		defer file.Close()

		err = template.Execute(file, properties)
		if err != nil {
			return err
		}
	}
	return nil
}
