// test Login
package main

import (
	"log"
	"os"

	"github.com/letrongvu/blog/internal/repository"
	"github.com/letrongvu/blog/internal/services"
)

func main() {
	db, err := repository.NewDB("./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userService := services.NewUserService(repository.NewUserRepository(db), repository.NewSessionRepository(db))
	if err := userService.Register(os.Args[1], os.Args[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("seeded user:", os.Args[1])
}
