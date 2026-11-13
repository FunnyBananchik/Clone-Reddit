package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"redditclone/internal/handler"
	"redditclone/internal/post"
	"redditclone/internal/user"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("set JWT_SECRET befor starting the server")
	}
	app := http.NewServeMux()
	rep := post.NewRepository()
	newUser := user.NewRepository()

	fs := http.FileServer(http.Dir("./static"))
	app.Handle("/static/", http.StripPrefix("/static/", fs))
	app.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/favicon.ico")
	})
	app.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/manifest.json")
	})

	hand := handler.NewPostHandler(rep)
	app.HandleFunc("GET /api/posts/", hand.GetPosts)
	//app.HandleFunc("/api/posts", hand.GetPosts)
	app.HandleFunc("GET /api/post/{postID}", hand.GetPostId)
	app.HandleFunc("GET /api/posts/{category}", hand.GetPostCategory)
	auth := handler.NewAuthHandler(newUser, secret)
	app.HandleFunc("POST /api/register", auth.Register)
	app.HandleFunc("POST /api/login", auth.Login)
	authMiddlewear := auth.AuthMiddlewear(http.HandlerFunc(hand.CreatePost))
	app.Handle("POST /api/posts", authMiddlewear)
	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/index.html")
	})
	err := http.ListenAndServe(":8080", app)
	if err != nil {
		fmt.Println("Не удалось запустить сервер")
	}
}
