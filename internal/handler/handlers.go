package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"redditclone/internal/post"
	"redditclone/internal/user"
	"redditclone/internal/utils"
	"regexp"
	"sort"
	"time"
	"unicode/utf8"
)

type PostsInterface interface {
	GetPosts() []post.Post
	GetPostId(string) (post.Post, error)
	GetPostCategory(string) []post.Post
	CreatePost(post.Post) post.Post
	GetPostsByUser(string) []post.Post
}

type UsersInterface interface {
	CheckExists(string) bool
}

type PostHandler struct {
	posts PostsInterface
	users UsersInterface
}

func NewPostHandler(posts PostsInterface, users UsersInterface) *PostHandler {
	return &PostHandler{posts: posts, users: users}
}

type PostRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	URL      string `json:"url,omitempty"`
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

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Path: %s \n", r.Method, r.URL.Path)
	u, ok := r.Context().Value(UserContextKey).(user.User)
	if !ok {
		errorJSON(w, 401, "user unauthorized")
		return
	}
	author := post.User{
		Id:   u.Id,
		Name: u.Username,
	}
	var req *PostRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := decoder.Decode(&req); err != nil || req == nil {
		errorJSON(w, 400, "cannot decode payload")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		errorJSON(w, 400, "cannot decode payload")
		return
	}
	if err := CheckValid(`^\S(?:[\s\S]*\S)?$`, req.Title, "title", 1, 100); err != nil {
		errorValidation(w, 422, "title", "bad title")
		return
	}
	if req.Category == "" {
		errorValidation(w, 422, "category", "bad category")
		return
	}
	if req.Type != "text" && req.Type != "link" {
		errorValidation(w, 422, "type", "bad type")
		return
	}
	switch req.Type {
	case "text":
		if utf8.RuneCountInString(req.Text) < 4 {
			errorValidation(w, 422, "text", "bad text")
			return
		}
	case "link":
		if _, err := url.ParseRequestURI(req.URL); err != nil {
			errorValidation(w, 422, "link", "bad link")
			return
		}
	}
	newPost := post.Post{
		Id:       utils.GenerateID(),
		Title:    req.Title,
		Category: req.Category,
		Type:     req.Type,
		Text:     req.Text,
		URL:      req.URL,
		Author:   author,
		Score:    1,
		Votes: []post.Vote{
			{
				User: author.Id,
				Vote: 1,
			},
		},
		Comments:         []post.Comment{},
		Created:          time.Now(),
		Views:            0,
		UpvotePercentage: 100,
	}
	p := h.posts.CreatePost(newPost)
	writeJSON(w, 201, p)
}

func (h *PostHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if validUsername, _ := regexp.MatchString(`^[A-Za-z0-9_-]{1,32}$`, username); !validUsername {
		errorValidation(w, 400, "username", "bad username")
		return
	}
	if !h.users.CheckExists(username) {
		errorJSON(w, 404, "user not found")
		return
	}
	posts := h.posts.GetPostsByUser(username)
	if len(posts) > 0 {
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Created.After(posts[j].Created)
		})
	}
	writeJSON(w, 200, posts)
}

func CheckValid(pattern, param, nameParam string, min, max int) error {
	validParam, _ := regexp.MatchString(pattern, param)
	if !validParam {
		return fmt.Errorf("bad parameter")
	}
	if utf8.RuneCountInString(param) < min || utf8.RuneCountInString(param) > max {
		return fmt.Errorf("bad parameter")
	}
	return nil
}
