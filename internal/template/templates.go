package template

import (
	"embed"
	"strings"
	"text/template"
)

//go:embed *.tmpl
var templateFS embed.FS

// Get returns the parsed template by name (e.g., "agents.md").
func Get(name string) (*template.Template, error) {
	data, err := templateFS.ReadFile(name + ".tmpl")
	if err != nil {
		return nil, err
	}
	return template.New(name).Parse(string(data))
}

// Render renders a template by name with the given data and returns the result.
func Render(name string, data any) (string, error) {
	tmpl, err := Get(name)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MustRender renders a template or panics.
func MustRender(name string, data any) string {
	s, err := Render(name, data)
	if err != nil {
		panic(err)
	}
	return s
}

// Raw returns the raw template content without rendering.
func Raw(name string) (string, error) {
	data, err := templateFS.ReadFile(name + ".tmpl")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
