package template

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

const (
	// include all files in the directory
	includeAllInDirectory = "./"
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

		fileName := strings.TrimPrefix(strings.TrimPrefix(fpath, source.BasePath), string(os.PathSeparator))
		isFound := slices.IndexFunc(source.Files, func(indexEntry string) bool {
			// FIXME: this is a hack, we need to find a better way to handle this
			if indexEntry == includeAllInDirectory {
				// include all files in the directory
				return true
			}
			return indexEntry == fileName
		})
		if isFound == -1 {
			// not a file that is part of the template
			return nil
		}

		tmpl, err := parseFile(fpath)
		if err != nil {
			return err
		}

		template.Files[fileName] = tmpl
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
	for fileName, tmpl := range a.Files {
		baseFileName := filepath.Base(fileName)
		if len(fileName) > len(baseFileName) {
			// create the directory structure
			err := os.MkdirAll(path.Join(originPath, strings.TrimSuffix(fileName, baseFileName)), 0775)
			if err != nil {
				return err
			}
		}

		// create the file and render the template
		file, err := os.OpenFile(path.Join(originPath, fileName), os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
		defer file.Close()

		err = tmpl.Execute(file, properties)
		if err != nil {
			return fmt.Errorf("failed to render template file %s: %w", fileName, err)
		}
	}
	return nil
}
