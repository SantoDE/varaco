package render

import (
	"text/template"
	"github.com/SantoDE/varaco/types"
	"os"
	"fmt"
)

type configData struct {
	Containers []types.ContainerData
	Port string
	AdditionalVCL string
}

const tmpl = `
vcl 4.0;

import directors;    # load the directors

{{range .Containers}}
backend server{{ .Id }} {
    .host = "{{ .Ip }}";
	.port = "{{ .Port }}";
}
{{ end }}

sub vcl_init {
    new bar = directors.round_robin();
	{{range .Containers}}
	bar.add_backend(server{{ .Id }});
	{{ end }}
}

sub vcl_recv {
    set req.backend_hint = bar.backend();
}

{{ if .AdditionalVCL }} 
	include "{{ .AdditionalVCL }}";
{{end}}
`

func RenderConfig(path string, containers []types.ContainerData, additionalVCL string, port string) (error){

	f, err := os.Create(path)

	if err != nil {
		fmt.Printf("Error Creating File for Config %s \n", err.Error())
		return err
	}

	for k := range containers {
		container := &containers[k]
		container.Port = port
	}

	data := new(configData)
	data.Containers = containers
	data.AdditionalVCL = additionalVCL

	t, err := template.New("varnish.vcl").Parse(tmpl)
	if err != nil {
		fmt.Printf("Error Creating Template %s \n", err.Error())
		return err
	}
	if err = t.Execute(f, data); err != nil {
		fmt.Printf("Error Rendering Template %s \n", err.Error())
		return err
	}

	return nil
}