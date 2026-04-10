package main

import (
	"log"
	"net/http"

	"github.com/Bojidarist/linkor/internal/config"
	"github.com/Bojidarist/linkor/internal/database"
	"github.com/Bojidarist/linkor/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	handler := server.New(cfg, db)

	addr := ":" + cfg.Port
	log.Printf("Linkor starting on %s", addr)
	log.Printf("Admin panel: http://localhost%s/admin/management?key=<YOUR_KEY>", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
