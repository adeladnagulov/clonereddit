package main

import (
	"cloneRaddit/handlers"
	"cloneRaddit/middleware"
	"cloneRaddit/post"
	"cloneRaddit/user"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	userRepo := user.CreateMemoryRepo()
	userHandle := handlers.NewUserHandler(userRepo)

	postRepo := post.NewPostMemoryRepo()
	postHandle := handlers.NewPostHandle(postRepo)

	r := mux.NewRouter()

	webHtmlHandler := http.FileServer(http.Dir("./web"))
	r.Handle("/", webHtmlHandler)

	staticHandle := http.FileServer(http.Dir("./web"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandle))

	r.HandleFunc("/api/register", userHandle.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandle.Login).Methods("POST")
	authPostHandle := middleware.Auth(http.HandlerFunc(postHandle.AddPost))
	r.Handle("/api/posts", authPostHandle).Methods("POST")

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
