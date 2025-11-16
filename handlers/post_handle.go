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
		http.Error(w, "json Marifk error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
