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
	VotesPost        []*Vote               `json:"votes"`
	Comments         []*Comment            `json:"comments"`
	Created          time.Time             `json:"created"`
	Views            int                   `json:"views"`
	Type             string                `json:"type"`
	Text             string                `json:"text,omitempty"`
	URL              string                `json:"url,omitempty"`
	UpvotePercentage int                   `json:"upvotePercentage"`
}

type Vote struct {
	UserID string `json:"user"`
	Vote   int    `json:"vote"`
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
	DeletePost(post *Post, autor middleware.UserClaims) error
	PostsByUser(username string) ([]*Post, error)
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

func (r *MemoryRepo) PostsByUser(username string) ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rez := []*Post{}
	for _, p := range r.Posts {
		if p.Author.UserName == username {
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
		VotesPost:        make([]*Vote, 0),
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

func (r *MemoryRepo) DeletePost(post *Post, autor middleware.UserClaims) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, v := range r.Posts {
		if v.ID == post.ID && post.Author.ID == autor.ID {
			r.Posts = append(r.Posts[:i], r.Posts[i+1:]...)
			return nil
		}
	}
	return errors.New("not faund or access is restricted")
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

func (p *Post) DeleteComment(commentId string, user middleware.UserClaims) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Author.ID != user.ID {
		return errors.New("access is restricted")
	}
	for i, v := range p.Comments {
		if v.ID == commentId {
			p.Comments = append(p.Comments[:i], p.Comments[i+1:]...)
			return nil
		}
	}
	return errors.New("comment not found")
}

func (p *Post) Vote(user string, voteValue int) {
	voteIndx := -1
	for i, v := range p.VotesPost {
		if v.UserID == user {
			voteIndx = i
			break
		}
	}

	if voteIndx != -1 {
		vote := p.VotesPost[voteIndx]
		p.Score -= vote.Vote
		if voteValue == 0 {
			p.VotesPost = append(p.VotesPost[:voteIndx], p.VotesPost[voteIndx+1:]...)
		} else {
			p.Score += voteValue
			vote.Vote = voteValue
		}
	} else if voteValue != 0 {
		p.Score += voteValue
		p.VotesPost = append(p.VotesPost, &Vote{UserID: user, Vote: voteValue})
	}
}

func (p *Post) CalculateUpvotePercentage() {
	p.mu.Lock()
	defer p.mu.Unlock()
	upvotes := 0
	downvotes := 0
	for _, v := range p.VotesPost {
		switch v.Vote {
		case 1:
			upvotes++
		case -1:
			downvotes++
		}
	}
	totalVotes := upvotes + downvotes
	if totalVotes != 0 {
		p.UpvotePercentage = (upvotes * 100) / totalVotes
	} else {
		p.UpvotePercentage = 0
	}

}
