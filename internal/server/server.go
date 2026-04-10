package server

import (
	"database/sql"
	"io/fs"
	"net/http"

	"github.com/Bojidarist/linkor/internal/config"
	"github.com/Bojidarist/linkor/internal/handlers"
	"github.com/Bojidarist/linkor/internal/repository"
	"github.com/Bojidarist/linkor/internal/services"
	"github.com/Bojidarist/linkor/web"
)

func New(cfg *config.Config, db *sql.DB) http.Handler {
	linkRepo := repository.NewLinkRepository(db)
	linkService := services.NewLinkService(linkRepo)
	adminHandler := handlers.NewAdminHandler(linkService)
	redirectHandler := handlers.NewRedirectHandler(linkService)
	authMiddleware := handlers.AdminKeyMiddleware(cfg)

	mux := http.NewServeMux()

	staticFS, _ := fs.Sub(web.Assets, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.Handle("GET /admin/management", authMiddleware(http.HandlerFunc(adminHandler.ServePage)))
	mux.Handle("GET /admin/api/links", authMiddleware(http.HandlerFunc(adminHandler.ListLinks)))
	mux.Handle("POST /admin/api/links", authMiddleware(http.HandlerFunc(adminHandler.CreateLink)))
	mux.Handle("PUT /admin/api/links/{id}", authMiddleware(http.HandlerFunc(adminHandler.UpdateLink)))
	mux.Handle("DELETE /admin/api/links/{id}", authMiddleware(http.HandlerFunc(adminHandler.DeleteLink)))

	mux.Handle("GET /{shortURL}", redirectHandler)

	return mux
}
