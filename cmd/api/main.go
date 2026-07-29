package main

import (
	"log"
	"os"

	"github.com/letrongvu/blog/internal/view"
	"github.com/letrongvu/blog/web"
)

func main() {
	if err := view.Init(web.Templates); err != nil {
		log.Fatal(err)
	}

	if err := view.Render(os.Stdout, "home.html", "world"); err != nil {
		log.Fatal(err)
	}
}