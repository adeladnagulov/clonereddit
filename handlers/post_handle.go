package handlers

import (
	"cloneRaddit/middleware"
	"cloneRaddit/post"
	"encoding/json"
	"net/http"
)

type PostHandle struct {
	Repo post.Repo
}

func NewPostHandle(repo post.Repo) *PostHandle {
	return &PostHandle{
		Repo: repo,
	}
}

func (h *PostHandle) AddPost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		Category string `json:"category"`
		Type     string `json:"type"`
		Text     string `json:"text"`
		URL      string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Невалидный json", http.StatusBadRequest)
		return
	}

	u, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	newPost := h.Repo.CreateNewPost(u.ID, u.UserName, req.Category, req.Title, req.Type, req.Text, req.URL)
	resp, err := json.Marshal(newPost)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *PostHandle) GetPosts(w http.ResponseWriter, r *http.Request) {
	list, err := h.Repo.GetAllPosts()
	if err != nil {
		http.Error(w, "Get posts error"+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *PostHandle) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("CATEGORY_NAME")
	list, err := h.Repo.PostsByCategory(category)
	if err != nil {
		http.Error(w, "CategoryPosts error"+err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *PostHandle) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("USER_LOGIN")
	list, err := h.Repo.PostsByUser(username)
	if err != nil {
		http.Error(w, "CategoryPosts error"+err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *PostHandle) GetPostsByID(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("POST_ID")
	post, err := h.Repo.PostByID(postId)
	if err != nil {
		http.Error(w, "find PostByID error", http.StatusNotFound)
		return
	}
	post.Views++

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *PostHandle) DeletePost(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("POST_ID")
	post, err := h.Repo.PostByID(postId)
	if err != nil {
		http.Error(w, "find PostByID error", http.StatusNotFound)
		return
	}
	u, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}
	err = h.Repo.DeletePost(post, u)
	if err != nil {
		http.Error(w, "DeletePost error:"+err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *PostHandle) AddComment(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("POST_ID")
	post, err := h.Repo.PostByID(postId)
	if err != nil {
		http.Error(w, "find PostByID error", http.StatusNotFound)
		return
	}

	u, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	var body struct {
		Comment string `json:"comment"`
	}
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error": "bad request"}`, http.StatusBadRequest)
		return
	}

	post.AddComment(u, body.Comment)

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *PostHandle) DeleteComment(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("POST_ID")
	commentId := r.PathValue("COMMENT_ID")
	post, err := h.Repo.PostByID(postId)
	if err != nil {
		http.Error(w, "find PostByID error", http.StatusNotFound)
		return
	}

	u, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}
	err = post.DeleteComment(commentId, u)
	if err != nil {
		if err.Error() == "access is restricted" {
			http.Error(w, "access is restricted", http.StatusForbidden)
			return
		} else {
			http.Error(w, "comment not found", http.StatusNotFound)
			return
		}
	}

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *PostHandle) Upvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, 1)
}

func (h *PostHandle) Downvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, -1)
}

func (h *PostHandle) Unvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, 0)
}

func (h *PostHandle) vote(w http.ResponseWriter, r *http.Request, voteValue int) {
	postId := r.PathValue("POST_ID")
	u, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}
	post, err := h.Repo.PostByID(postId)
	if err != nil {
		http.Error(w, "find PostByID error", http.StatusNotFound)
		return
	}
	post.Vote(u.ID, voteValue)
	post.CalculateUpvotePercentage()

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "json Marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
