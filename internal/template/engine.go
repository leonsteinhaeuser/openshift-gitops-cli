package template

import (
	"os"
	"path"
	"text/template"

	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/menu"
)

type Template struct {
	Path             string
	TemplateManifest TemplateManifest
}

type TemplateCarrier struct {
	TemplateName string
	FileName     string
	Template     *template.Template
}

// Render renders the template with the given carrier
func (t Template) Render(basePath string, ccc menu.CarrierCreateCluster) error {
	files, err := loadAsTemplate(t)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = renderTemplate(basePath, ccc, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadAsTemplate loads the template files as a template
func loadAsTemplate(t Template) ([]TemplateCarrier, error) {
	files := []TemplateCarrier{}
	for _, file := range t.TemplateManifest.Files {
		fpath := path.Join(t.Path, file)
		// FIXME: this is not efficient
		bts, err := os.ReadFile(fpath)
		if err != nil {
			return nil, err
		}
		// parse the template
		tpl, err := template.New("template").Parse(string(bts))
		if err != nil {
			return nil, err
		}
		files = append(files, TemplateCarrier{
			TemplateName: t.TemplateManifest.Name,
			FileName:     file,
			Template:     tpl,
		})
	}
	return files, nil
}

// renderTemplate renders the template with the given carrier and writes it to the file system
func renderTemplate(basePath string, ccc menu.CarrierCreateCluster, t TemplateCarrier) error {
	dpath := path.Join(basePath, ccc.Environment, ccc.Stage, ccc.ClusterName, t.TemplateName)
	err := os.MkdirAll(dpath, 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(dpath, t.FileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Template.Execute(file, ccc)
	if err != nil {
		return err
	}
	return nil
}
