package view

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/ninedraft/daily-bacon/internal/models"
)

//go:embed *.template
var fsys embed.FS

var parsed = template.Must(template.New("").
	Funcs(funcs).
	ParseFS(fsys, "*.template"))

var airQuaility = mustFind("AirQuality")

func AirQuality(dst io.Writer, params models.AirQualityResponse) error {
	return airQuaility.Execute(dst, params)
}

func mustFind(name string) *template.Template {
	t := parsed.Lookup(name)
	if t == nil {
		panic(fmt.Sprintf("unable to find template %q. Available templates: %s", name, parsed.DefinedTemplates()))
	}

	return t
}
