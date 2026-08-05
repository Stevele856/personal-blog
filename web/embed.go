package web

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed templates
var templatesDir embed.FS

//go:embed static
var staticDir embed.FS

var Templates fs.FS = resolveFS(templatesDir, "templates", "web/templates")
var Static fs.FS = resolveFS(staticDir, "static", "web/static")

func resolveFS(embedded embed.FS, subdir, diskPath string) fs.FS{
	if os.Getenv("APP_ENV") == "dev"{
		return os.DirFS(diskPath)
	}
	return must(fs.Sub(embedded,subdir))
}

func must(fsys fs.FS, err error) fs.FS{
	if err != nil {
		panic(err)
	}
	return fsys
}