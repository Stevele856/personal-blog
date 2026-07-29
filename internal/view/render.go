package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path"
)

var funcMap = template.FuncMap{}
var pages map[string]*template.Template

func Init(fsys fs.FS) error {
	base, err := template.New("layout").Funcs(funcMap).ParseFS(fsys,
		"layouts/*.html",
		"partials/*.html",
	)
	if err != nil {
		return fmt.Errorf("view: parsing base: %w", err)
	}

	pageFiles, err := fs.Glob(fsys, "pages/*.html")

	if err != nil {
		return fmt.Errorf("view: globbing pages: %w", err)
	}

	adminFiles, err := fs.Glob(fsys, "pages/admin/*.html")
	if err != nil {
		return fmt.Errorf("view: globbing admin pages: %w", err)
	}
	pageFiles = append(pageFiles, adminFiles...)

	pages = make(map[string]*template.Template, len(pageFiles))

	for _, f := range pageFiles {
		clone, err := base.Clone()
		if err != nil {
			return fmt.Errorf("view: cloning base for %s: %w", f, err)
		}

		clone, err = clone.ParseFS(fsys, f)
		if err != nil {
			return fmt.Errorf("view: parsing %s: %w", f, err)
		}
		pages[path.Base(f)] = clone
	}

	if len(pages) == 0 {
		return fmt.Errorf("view: no page templates found")
	}
	return nil
}

func Render(w io.Writer, page string, data any) error {
	t, ok := pages[page]
	if !ok {
		return fmt.Errorf("view: page %q not registered", page)
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout", data); err != nil {
		return fmt.Errorf("view: rendering %s: %w", page, err)
	}
	_, err := buf.WriteTo(w)
	return  err
}
