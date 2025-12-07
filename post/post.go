package post

import (
	"cloneRaddit/middleware"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID               string                `json:"id"`
	Title            string                `json:"title"`
	Author           middleware.UserClaims `json:"author"`
	Category         string                `json:"category"`
	Score            int                   `json:"score"`
	VotesPost        []Votes               `json:"votes"`
	Comments         []CommentsPost        `json:"comments"`
	Created          time.Time             `json:"created"`
	Views            int                   `json:"views"`
	Type             string                `json:"type"`
	Text             string                `json:"text,omitempty"`
	URL              string                `json:"url,omitempty"`
	UpvotePercentage int                   `json:"upvotePercentage"`
}

type Votes struct { // не доработан
	User string
	Vote int
}

type CommentsPost struct { //заглушка

}

type Repo interface {
	CreateNewPost(userId, username, category, title, postType, text, url string) *Post
	GetAllPosts() ([]*Post, error)
	CategoryPosts(category string) ([]*Post, error)
}

type MemoryRepo struct {
	mu    sync.RWMutex
	Posts []*Post
}

func NewPostMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		Posts: make([]*Post, 0),
	}
}

func (r *MemoryRepo) GetAllPosts() ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Posts, nil
}

func (r *MemoryRepo) CategoryPosts(category string) ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rez := []*Post{}
	for _, p := range r.Posts {
		if p.Category == category {
			rez = append(rez, p)
		}
	}
	return rez, nil
}

func (r *MemoryRepo) CreateNewPost(userId, username, category, title, postType, text, url string) *Post {
	r.mu.Lock()
	defer r.mu.Unlock()

	post := Post{
		Author: middleware.UserClaims{
			ID:       userId,
			UserName: username,
		},
		Category:         category,
		Created:          time.Now(),
		ID:               uuid.NewString(),
		Score:            1,
		Title:            title,
		Type:             postType,
		UpvotePercentage: 100,
		Views:            0,
		VotesPost:        make([]Votes, 0),
		Comments:         make([]CommentsPost, 0),
	}
	if postType == "text" {
		post.Text = text
	} else {
		post.URL = url
	}
	r.Posts = append(r.Posts, &post)
	return &post
}
