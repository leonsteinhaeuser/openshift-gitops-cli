package template

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Template struct {
	// Path defines the location to the manifest.yaml file
	Path string
	// TemplateManifest is the content of the manifest.yaml file
	TemplateManifest TemplateManifest
}

type TemplateCarrier struct {
	TemplateName string
	FileName     string
	Template     *template.Template
}

type TemplateData struct {
	Environment string
	Stage       string
	ClusterName string
	Addons      map[string]AddonData
	Properties  map[string]string
}

type AddonData struct {
	Annotations map[string]string
	Properties  map[string]any
}

// Render renders the template with the given carrier
func (t Template) Render(basePath string, td TemplateData) error {
	files, err := t.loadAsTemplate()
	if err != nil {
		return err
	}
	for _, file := range files {
		err = renderTemplate(basePath, td, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadAsTemplate loads the template files as a template
func (t Template) loadAsTemplate() ([]TemplateCarrier, error) {
	files := []TemplateCarrier{}
	for _, file := range t.TemplateManifest.Files {
		fpath := path.Join(t.Path, file)

		finfo, err := os.Stat(fpath)
		if err != nil {
			return nil, err
		}

		if !finfo.IsDir() {
			// is a file
			tmpl, err := parseFile(fpath)
			if err != nil {
				return nil, err
			}
			tc := TemplateCarrier{
				TemplateName: t.TemplateManifest.Name,
				FileName:     file,
				Template:     tmpl,
			}
			files = append(files, tc)
			continue
		}

		// if a directory
		// load all files in the directory
		err = filepath.WalkDir(fpath, func(fpath string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				// we don't care about directories
				return nil
			}
			tmpl, err := parseFile(fpath)
			if err != nil {
				return err
			}
			tc := TemplateCarrier{
				TemplateName: t.TemplateManifest.Name,
				FileName:     strings.TrimPrefix(fpath, t.Path),
				Template:     tmpl,
			}
			files = append(files, tc)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func parseFile(fpath string) (*template.Template, error) {
	bts, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	// parse the template
	tpl, err := template.New("template").Funcs(funcMap()).Parse(string(bts))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

// renderTemplate renders the template with the given carrier and writes it to the file system
func renderTemplate(basePath string, td TemplateData, t TemplateCarrier) error {
	dpath := path.Join(basePath, td.Environment, td.Stage, td.ClusterName, t.TemplateName)
	if fd := path.Dir(t.FileName); fd != "." {
		t.FileName = strings.TrimPrefix(t.FileName, fd)
		dpath = path.Join(dpath, fd)
	}
	err := os.MkdirAll(dpath, 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(dpath, t.FileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Template.Execute(file, td)
	if err != nil {
		return err
	}
	return nil
}
