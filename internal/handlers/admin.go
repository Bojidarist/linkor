package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Bojidarist/linkor/internal/models"
	"github.com/Bojidarist/linkor/internal/services"
	"github.com/Bojidarist/linkor/web"
)

type AdminHandler struct {
	linkService *services.LinkService
}

func NewAdminHandler(linkService *services.LinkService) *AdminHandler {
	return &AdminHandler{linkService: linkService}
}

func (h *AdminHandler) ServePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := web.Assets.ReadFile("templates/admin.html")
	if err != nil {
		log.Printf("reading admin template: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(tmpl)
}

func (h *AdminHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.linkService.List()
	if err != nil {
		writeJSONError(w, "failed to list links", http.StatusInternalServerError)
		log.Printf("listing links: %v", err)
		return
	}
	if links == nil {
		links = []models.Link{}
	}
	writeJSON(w, links, http.StatusOK)
}

func (h *AdminHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLinkRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.linkService.Create(req)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, link, http.StatusCreated)
}

func (h *AdminHandler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, "invalid link id", http.StatusBadRequest)
		return
	}

	var req models.UpdateLinkRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.linkService.Update(id, req)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, link, http.StatusOK)
}

func (h *AdminHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, "invalid link id", http.StatusBadRequest)
		return
	}

	if err := h.linkService.Delete(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encoding JSON response: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, map[string]string{"error": message}, status)
}

func readJSON(r *http.Request, v any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
