package main

import (
	"log"
	"net/http"
	"os"

	"github.com/letrongvu/blog/internal/handler"
	"github.com/letrongvu/blog/internal/middleware"
	"github.com/letrongvu/blog/internal/repository"
	"github.com/letrongvu/blog/internal/services"
	"github.com/letrongvu/blog/internal/view"
	"github.com/letrongvu/blog/web"
)

func main() {
	if err := view.Init(web.Templates); err != nil {
		log.Fatal(err)
	}

	// fallback
	dbPath := os.Getenv("DB_PATH")
	if dbPath == ""{
		dbPath = "./blog.db"
	}

	// Call NewDB
	db, err := repository.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	postRepo := repository.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postHandler := handler.New(postService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", postHandler.Home)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.Recover(mux)))

	// if err := view.Render(os.Stdout, "home.html", "world"); err != nil {
	// 	log.Fatal(err)
	// }
}