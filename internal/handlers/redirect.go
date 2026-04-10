package handlers

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Bojidarist/linkor/internal/services"
)

type RedirectHandler struct {
	linkService *services.LinkService
}

func NewRedirectHandler(linkService *services.LinkService) *RedirectHandler {
	return &RedirectHandler{linkService: linkService}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	if shortURL == "" {
		http.NotFound(w, r)
		return
	}

	clientIP := extractClientIP(r)
	targetURL, err := h.linkService.HandleRedirect(shortURL, clientIP)
	if err != nil {
		if strings.Contains(err.Error(), "link not found") {
			http.NotFound(w, r)
			return
		}
		log.Printf("redirect error for %q: %v", shortURL, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, targetURL, http.StatusFound)
}

func extractClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.SplitN(forwarded, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
