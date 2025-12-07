package main

import (
	"cloneRaddit/handlers"
	"cloneRaddit/middleware"
	"cloneRaddit/post"
	"cloneRaddit/user"
	"fmt"
	"log"
	"net/http"
)

func main() {
	userRepo := user.CreateMemoryRepo()
	userHandle := handlers.NewUserHandler(userRepo)

	postRepo := post.NewPostMemoryRepo()
	postHandle := handlers.NewPostHandle(postRepo)

	r := http.NewServeMux()
	r.HandleFunc("POST /api/register", userHandle.Register)
	r.HandleFunc("POST /api/login", userHandle.Login)
	r.HandleFunc("GET /api/posts", postHandle.GetPosts)
	r.HandleFunc("GET /api/posts/{CATEGORY_NAME}", postHandle.GetPostsByCategory)
	r.HandleFunc("GET /api/post/{POST_ID}", postHandle.GetPostsByID)

	authR := http.NewServeMux()
	authR.HandleFunc("POST /api/posts", postHandle.AddPost)
	authR.HandleFunc("POST /api/post/{POST_ID}", postHandle.AddComment)
	authR.HandleFunc("DELETE /api/post/{POST_ID}/{COMMENT_ID}", postHandle.DeleteComment)

	r.Handle("/api/", middleware.Auth(authR))

	webHtmlHandler := http.FileServer(http.Dir("./web/html"))
	r.Handle("/", webHtmlHandler)
	staticHandle := http.FileServer(http.Dir("./web"))
	r.Handle("/static/", http.StripPrefix("/static/", staticHandle))

	fmt.Println("starting server at :8080")

	handler := middleware.StripTrailingSlash(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
