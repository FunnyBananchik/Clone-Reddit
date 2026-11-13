package post

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Repository struct {
	posts []Post
	mu    sync.RWMutex
}

func NewRepository() *Repository {
	posts := []Post{
		{
			Id:               "507f1f77bcf86cd799439014",
			Title:            "Learning Go",
			Category:         "programming",
			Type:             "text",
			Text:             "What are your favorite resources for learning Go?",
			Author:           User{Id: "507f1f77bcf86cd799439011", Name: "alice"},
			Score:            1,
			Votes:            []Vote{{User: "507f1f77bcf86cd799439011", Vote: 1}},
			Comments:         []Comment{},
			Created:          time.Date(2026, time.November, 2, 12, 10, 0, 0, time.UTC),
			Views:            0,
			UpvotePercentage: 100,
		},
		{
			Id:               "507f1f77bcf86cd799439012",
			Title:            "The Go programming language",
			Category:         "programming",
			Type:             "link",
			URL:              "https://go.dev/",
			Author:           User{Id: "507f1f77bcf86cd799439011", Name: "alice"},
			Score:            5,
			Votes:            []Vote{{User: "507f1f77bcf86cd799439011", Vote: 1}},
			Comments:         []Comment{},
			Created:          time.Date(2026, time.November, 2, 12, 20, 0, 0, time.UTC),
			Views:            0,
			UpvotePercentage: 100,
		},
	}
	return &Repository{posts: posts}
}

func (list *Repository) GetPosts() []Post {
	list.mu.RLock()
	defer list.mu.RUnlock()
	if list == nil {
		return []Post{}
	}
	return list.posts
}

func (list *Repository) GetPostId(targetId string) (Post, error) {
	list.mu.Lock()
	defer list.mu.Unlock()
	for _, post := range list.posts {
		if post.Id == targetId {
			post.Views++
			return post, nil
		}
	}
	log.Print("Post not found")
	return Post{}, errors.New("Post not found")
}

func (list *Repository) GetPostCategory(targetCat string) []Post {
	list.mu.Lock()
	defer list.mu.Unlock()
	res := make([]Post, 0)
	for _, post := range list.posts {
		if post.Category == targetCat {
			res = append(res, post)
		}
	}
	return res
}

func (list *Repository) CreatePost(p Post) Post {
	list.mu.Lock()
	defer list.mu.Unlock()
	list.posts = append(list.posts, p)
	return p
}

func (list *Repository) GetPostsByUser(username string) []Post {
	list.mu.RLock()
	defer list.mu.RUnlock()
	posts := make([]Post, 0)
	for _, val := range list.posts {
		if val.Author.Name == username {
			posts = append(posts, val)
		}
	}
	return posts
}
