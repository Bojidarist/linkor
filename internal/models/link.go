package models

import "time"

type Link struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ShortURL     string    `json:"short_url"`
	TargetURL    string    `json:"target_url"`
	Clicks       int64     `json:"clicks"`
	UniqueClicks int64     `json:"unique_clicks"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateLinkRequest struct {
	Name      string `json:"name"`
	ShortURL  string `json:"short_url"`
	TargetURL string `json:"target_url"`
}

type UpdateLinkRequest struct {
	Name      string `json:"name"`
	ShortURL  string `json:"short_url"`
	TargetURL string `json:"target_url"`
}
