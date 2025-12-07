package post

import (
	"cloneRaddit/middleware"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	mu               sync.RWMutex
	ID               string                `json:"id"`
	Title            string                `json:"title"`
	Author           middleware.UserClaims `json:"author"`
	Category         string                `json:"category"`
	Score            int                   `json:"score"`
	VotesPost        []*Votes              `json:"votes"`
	Comments         []*Comment            `json:"comments"`
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

type Comment struct {
	ID      string                `json:"id"`
	Author  middleware.UserClaims `json:"author"`
	Body    string                `json:"body"`
	Created time.Time             `json:"created"`
}

type Repo interface {
	CreateNewPost(userId, username, category, title, postType, text, url string) *Post
	GetAllPosts() ([]*Post, error)
	PostsByCategory(category string) ([]*Post, error)
	PostByID(id string) (*Post, error)
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

func (r *MemoryRepo) PostsByCategory(category string) ([]*Post, error) {
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

func (r *MemoryRepo) PostByID(id string) (*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Posts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("Not found user")
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
		VotesPost:        make([]*Votes, 0),
		Comments:         make([]*Comment, 0),
	}
	if postType == "text" {
		post.Text = text
	} else {
		post.URL = url
	}
	r.Posts = append(r.Posts, &post)
	return &post
}

func (p *Post) AddComment(autor middleware.UserClaims, body string) {
	comment := &Comment{
		ID:      uuid.NewString(),
		Author:  autor,
		Body:    body,
		Created: time.Now(),
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Comments = append(p.Comments, comment)
}
