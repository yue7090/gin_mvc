package middleware

import (
	"net/http"
	"path"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin/render"
	"gopkg.in/ini.v1"
	"os"
	"fmt"
)

type RenderOptions struct {
	TemplateDir string
	ContentType string
}

// Pongo2Render is a custom Gin template renderer using Pongo2.
type Pongo2Render struct {
	Options  *RenderOptions
	Template *pongo2.Template
	Context  pongo2.Context
}

func New(options RenderOptions) *Pongo2Render {
	return &Pongo2Render{
		Options: &options,
	}
}


func DefaultTemplateDir() *Pongo2Render {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
	}
	template_dir := cfg.Section("web_dir").Key("template_dir").String()
	if template_dir == "" {
		fmt.Println("template_dir setting can not be null")
        os.Exit(1)
	}
	return New(RenderOptions{
		TemplateDir: template_dir,
		ContentType: "text/html; charset=utf-8",
	})
}

func (p Pongo2Render) Instance(name string, data interface{}) render.Render {
	var template *pongo2.Template
	filename := path.Join(p.Options.TemplateDir, name)
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
	}
	app_mode := cfg.Section("app").Key("app_mode").String()
	// always read template files from disk if in debug mode, use cache otherwise.
	if app_mode == "debug" {
		template = pongo2.Must(pongo2.FromFile(filename))
	} else {
		template = pongo2.Must(pongo2.FromCache(filename))
	}

	return Pongo2Render{
		Template: template,
		Context:  data.(pongo2.Context),
		Options:  p.Options,
	}
}

func (p Pongo2Render) Render(w http.ResponseWriter) error {
	p.WriteContentType(w)
	err := p.Template.ExecuteWriter(p.Context, w)
	return err
}

func (p Pongo2Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{p.Options.ContentType}
	}
}
