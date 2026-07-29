package web

import (
	"embed"
	"io/fs"
)

//go:embed templates
var templatesDir embed.FS
var Templates fs.FS = must(fs.Sub(templatesDir, "templates"))

//go:embed static
var staticDir embed.FS
var Static fs.FS = must(fs.Sub(staticDir, "static"))

func must(fsys fs.FS, err error) fs.FS{
	if err != nil {
		panic(err)
	}
	return fsys
}