package post

import "time"

type User struct {
	Id    string `json:"id"`
	Name  string `json:"username"`
	Admin bool   `json:"admin,omitempty"`
}

type Vote struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type Comment struct {
	Id      string    `json:"id"`
	Author  User      `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

type Post struct {
	Id               string    `json:"id"`
	Title            string    `json:"title"`
	Category         string    `json:"category"`
	Type             string    `json:"type"`
	Text             string    `json:"text,omitempty"`
	URL              string    `json:"url,omitempty"`
	Author           User      `json:"author"`
	Score            int       `json:"score"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"`
	Created          time.Time `json:"created"`
	Views            int       `json:"views"`
	UpvotePercentage int       `json:"upvotePercentage"`
}
