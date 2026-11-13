package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"redditclone/internal/post"
	"sort"
)

type PostsInterface interface {
	GetPosts() []post.Post
	GetPostId(string) (post.Post, error)
	GetPostCategory(string) []post.Post
}

type PostHandler struct {
	posts PostsInterface
}

func NewPostHandler(posts PostsInterface) *PostHandler {
	return &PostHandler{posts: posts}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Path: %s \n", r.Method, r.URL.Path)
	res := h.posts.GetPosts()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Score > res[j].Score
	})
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("error encode posts: %v\n", err)
	}
}

func (h *PostHandler) GetPostId(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Path: %s \n", r.Method, r.URL.Path)
	id := r.PathValue("postID")
	w.Header().Set("Content-Type", "application/json")
	targetPost, err := h.posts.GetPostId(id)
	if err != nil {
		w.WriteHeader(404)
		textError := fmt.Sprintf("post with id %q not found", id)
		log.Printf("error: %s\n", textError)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": textError}); err != nil {
			log.Printf("error encode posts: %v\n", err)
		}
		return
	}
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(targetPost); err != nil {
		log.Printf("error encode posts: %v\n", err)
	}
}

func (h *PostHandler) GetPostCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Path: %s \n", r.Method, r.URL.Path)
	category := r.PathValue("category")
	w.Header().Set("Content-Type", "application/json")
	targetPosts := h.posts.GetPostCategory(category)
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(targetPosts); err != nil {
		log.Printf("error encode posts: %v\n", err)
	}
}

/*func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Path: %s \n", r.Method, r.URL.Path)
	author, ok := r.Context().Value(UserContextKey).(post.User)
	if !ok {
		errorJSON(w, 401, "user unauthorized")
	}

}
*/
