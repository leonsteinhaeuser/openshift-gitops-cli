package template

import (
	"bytes"
	"compress/gzip"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

func funcMap() template.FuncMap {
	templateFuncMap := sprig.FuncMap()
	templateFuncMap["toYaml"] = toYAML
	templateFuncMap["gzip"] = gzipCompress
	templateFuncMap["gunzip"] = gzipDecompress
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
	buffer := new(bytes.Buffer)
	w := gzip.NewWriter(buffer)
	defer w.Close()
	_, err := w.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func gzipDecompress(data string) string {
	r, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
