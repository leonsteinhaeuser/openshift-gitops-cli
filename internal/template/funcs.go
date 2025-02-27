package template

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

func funcMap() template.FuncMap {
	templateFuncMap := sprig.FuncMap()
	templateFuncMap["toYaml"] = toYAML
	return templateFuncMap
}

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
