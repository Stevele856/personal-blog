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
	"github.com/letrongvu/blog/migration"
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

	// Open NewDB
	db, err := repository.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Call migrate after open DB, before create repo
	if err := repository.Migrate(db,migration.FS); err != nil {
		log.Fatal(err)
	}

	postRepo := repository.NewPostRepository(db)
	postService := services.NewPostService(postRepo)

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	userService := services.NewUserService(userRepo, sessionRepo)

	postHandler := handler.New(postService, userService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", postHandler.Home)
	mux.HandleFunc("GET /about", postHandler.About)
	mux.HandleFunc("GET /posts/{slug}", postHandler.Post)

	mux.HandleFunc("GET /login", postHandler.LoginForm)
	mux.HandleFunc("POST /login", postHandler.Login)
	mux.HandleFunc("POST /logout", postHandler.Logout)

	auth := middleware.RequireAuth(userService)
	mux.Handle("GET /admin", auth(http.HandlerFunc(postHandler.Dashboard)))
	mux.Handle("GET /admin/posts/new", auth(http.HandlerFunc(postHandler.NewPostForm)))
	mux.Handle("POST /admin/posts", auth(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("GET /admin/posts/{id}/edit", auth(http.HandlerFunc(postHandler.EditPostForm)))
	mux.Handle("POST /admin/posts/{id}", auth(http.HandlerFunc(postHandler.UpdatePost)))
	mux.Handle("POST /admin/posts/{id}/delete", auth(http.HandlerFunc(postHandler.DeletePost)))
	
	mux.HandleFunc("/", postHandler.NotFound)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.Recover(mux)))

	// if err := view.Render(os.Stdout, "home.html", "world"); err != nil {
	// 	log.Fatal(err)
	// }
}