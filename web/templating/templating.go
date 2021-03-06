package templating

import (
	"bones/config"
	"html/template"
	"io"
	"log"
)

var templates *template.Template

type TemplateRenderer interface {
	RenderTemplate(wr io.Writer, ctx TemplateContext) error
}

type TemplateContext interface {
	// Add an error that should be displayed
	// to the user
	AddError(err error)

	// Add a notice/message that should be
	// displayed to the user
	AddNotice(notice string)

	// Name of the (main) template that is being
	// rendered (not header or footer templates).
	Name() string
}

func NewTemplateRenderer() TemplateRenderer {
	if config.Env().IsProduction() {
		return &cachingTemplateRenderer{templates}
	}

	log.Println("Using reloading template renderer")

	return &reloadingTemplateRenderer{}
}

type cachingTemplateRenderer struct {
	templates *template.Template
}

func (r cachingTemplateRenderer) RenderTemplate(wr io.Writer, ctx TemplateContext) error {
	if r.templates == nil {
		var err error
		r.templates, err = template.ParseGlob("./templates/*.html")

		if err != nil {
			return err
		}
	}

	return renderWithLayout(r.templates, wr, ctx.Name(), ctx)
}

type reloadingTemplateRenderer struct{}

func (r reloadingTemplateRenderer) RenderTemplate(wr io.Writer, ctx TemplateContext) error {
	templates, err := template.ParseGlob("./templates/*.html")

	if err != nil {
		return err
	}

	return renderWithLayout(templates, wr, ctx.Name(), ctx)
}

func renderWithLayout(t *template.Template, wr io.Writer, name string, data interface{}) error {
	for _, templateName := range []string{"header.html", name, "footer.html"} {
		err := t.ExecuteTemplate(wr, templateName, data)

		if err != nil {
			return err
		}
	}

	return nil
}
