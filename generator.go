// +build ignore

package main

import (
	fun "functions"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var packageTemplate = template.Must(template.New("").
	Parse("// Code generated by go generate; DO NOT EDIT.\n" +
		"package main\n" +
		"\n" +
		"var hasgoTemplates = map[string]string{\n" +
		"{{ range $fn, $file := . }}" +
		"\t\"{{ $fn }}\": `{{ $file }}`,\n" +
		"{{ end }}" +
		"}\n"))

var domainTemplate = template.Must(template.New("").
	Parse("\n" +
		"const (\n ForNumbers = \"ForNumbers\"\nForStrings = \"ForStrings\"\n" +
		"ForStructs = \"ForStructs\"\n)\n" +
		"var funcDomains = map[string][]string{\n" +
		"{{ range $fn, $arr := . }}" +
		"\t\"{{ $fn }}\": []string{ {{ range $index, $dom := $arr }}" +
		" {{if $index}} ,{{end}} {{$dom}} {{end}} },\n" +
		"{{ end }}" +
		"}\n"))

// generate the templates that will be used for hasgo
func main() {
	data := map[string]string{}
	for k := range fun.Templates() {
		content, err := ioutil.ReadFile("./functions/" + k)
		// remove the package statement..
		parts := strings.Split(string(content), "\n")
		sanitized := strings.Join(parts[1:], "\n")
		if err != nil {
			panic(err)
		}
		data[k] = sanitized
	}
	f, err := os.Create("template.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	packageTemplate.Execute(f, data)
	domainTemplate.Execute(f, fun.Templates())
}
