package repository

import (
	"database/sql"
	"fmt"

	"github.com/Bojidarist/linkor/internal/models"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) GetAll() ([]models.Link, error) {
	rows, err := r.db.Query(
		`SELECT id, name, short_url, target_url, clicks, unique_clicks, created_at, updated_at
		 FROM links ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("querying links: %w", err)
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		if err := rows.Scan(&l.ID, &l.Name, &l.ShortURL, &l.TargetURL, &l.Clicks, &l.UniqueClicks, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning link: %w", err)
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

func (r *LinkRepository) GetByShortURL(shortURL string) (*models.Link, error) {
	var l models.Link
	err := r.db.QueryRow(
		`SELECT id, name, short_url, target_url, clicks, unique_clicks, created_at, updated_at
		 FROM links WHERE short_url = ?`, shortURL,
	).Scan(&l.ID, &l.Name, &l.ShortURL, &l.TargetURL, &l.Clicks, &l.UniqueClicks, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("querying link by short_url: %w", err)
	}
	return &l, nil
}

func (r *LinkRepository) GetByID(id int64) (*models.Link, error) {
	var l models.Link
	err := r.db.QueryRow(
		`SELECT id, name, short_url, target_url, clicks, unique_clicks, created_at, updated_at
		 FROM links WHERE id = ?`, id,
	).Scan(&l.ID, &l.Name, &l.ShortURL, &l.TargetURL, &l.Clicks, &l.UniqueClicks, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("querying link by id: %w", err)
	}
	return &l, nil
}

func (r *LinkRepository) Create(name, shortURL, targetURL string) (*models.Link, error) {
	result, err := r.db.Exec(
		`INSERT INTO links (name, short_url, target_url) VALUES (?, ?, ?)`,
		name, shortURL, targetURL,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last insert id: %w", err)
	}
	return r.GetByID(id)
}

func (r *LinkRepository) Update(id int64, name, shortURL, targetURL string) (*models.Link, error) {
	_, err := r.db.Exec(
		`UPDATE links SET name = ?, short_url = ?, target_url = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		name, shortURL, targetURL, id,
	)
	if err != nil {
		return nil, fmt.Errorf("updating link: %w", err)
	}
	return r.GetByID(id)
}

func (r *LinkRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM links WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting link: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("link with id %d not found", id)
	}
	return nil
}

func (r *LinkRepository) ShortURLExists(shortURL string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM links WHERE short_url = ? AND id != ?`,
		shortURL, excludeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking short_url existence: %w", err)
	}
	return count > 0, nil
}

func (r *LinkRepository) IncrementClicks(id int64) error {
	_, err := r.db.Exec(`UPDATE links SET clicks = clicks + 1 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("incrementing clicks: %w", err)
	}
	return nil
}

func (r *LinkRepository) IncrementUniqueClicks(id int64) error {
	_, err := r.db.Exec(`UPDATE links SET unique_clicks = unique_clicks + 1 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("incrementing unique clicks: %w", err)
	}
	return nil
}

func (r *LinkRepository) RecordClick(linkID int64, ipHash string) (isNew bool, err error) {
	result, err := r.db.Exec(
		`INSERT OR IGNORE INTO click_tracking (link_id, ip_hash) VALUES (?, ?)`,
		linkID, ipHash,
	)
	if err != nil {
		return false, fmt.Errorf("recording click: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("checking rows affected: %w", err)
	}
	return rows > 0, nil
}
