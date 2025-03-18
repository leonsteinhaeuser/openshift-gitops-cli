package template

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

func funcMap(tmpl *template.Template) template.FuncMap {
	templateFuncMap := sprig.FuncMap()
	templateFuncMap["toYaml"] = toYAML
	templateFuncMap["gzip"] = gzipCompress
	templateFuncMap["gunzip"] = gzipDecompress
	templateFuncMap["include"] = includeFun(tmpl, map[string]int{})
	templateFuncMap["joinPath"] = path.Join
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

func gzipCompress(data string) string {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write([]byte(data))
	if err != nil {
		panic(err)
	}

	// Close the gzip writer to flush any remaining data
	if err := gz.Close(); err != nil {
		panic(err)
	}

	return buf.String()
}

func gzipDecompress(data string) string {
	// Create a reader from the compressed data
	buf := bytes.NewReader([]byte(data))
	gz, err := gzip.NewReader(buf)
	if err != nil {
		panic(err)
	}
	defer gz.Close()

	// Read and decompress the data
	var out bytes.Buffer
	_, err = io.Copy(&out, gz)
	if err != nil {
		panic(err)
	}

	return out.String()
}

const recursionMaxNums = 1000

// include a define block
func includeFun(t *template.Template, includedNames map[string]int) func(string, interface{}) (string, error) {
	return func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", fmt.Errorf("unable to execute template: rendering template has a nested reference name: %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := t.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}
}
