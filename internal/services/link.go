package services

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/Bojidarist/linkor/internal/models"
	"github.com/Bojidarist/linkor/internal/repository"
)

const (
	shortURLLength  = 6
	shortURLCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetries      = 10
)

type LinkService struct {
	repo *repository.LinkRepository
}

func NewLinkService(repo *repository.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func (s *LinkService) List() ([]models.Link, error) {
	return s.repo.GetAll()
}

func (s *LinkService) Create(req models.CreateLinkRequest) (*models.Link, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.TargetURL) == "" {
		return nil, fmt.Errorf("target URL is required")
	}

	shortURL := strings.TrimSpace(req.ShortURL)
	if shortURL == "" {
		generated, err := s.generateUniqueShortURL()
		if err != nil {
			return nil, err
		}
		shortURL = generated
	} else {
		if err := validateShortURL(shortURL); err != nil {
			return nil, err
		}
		exists, err := s.repo.ShortURLExists(shortURL, 0)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("short URL %q is already in use", shortURL)
		}
	}

	return s.repo.Create(strings.TrimSpace(req.Name), shortURL, strings.TrimSpace(req.TargetURL))
}

func (s *LinkService) Update(id int64, req models.UpdateLinkRequest) (*models.Link, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.TargetURL) == "" {
		return nil, fmt.Errorf("target URL is required")
	}
	if strings.TrimSpace(req.ShortURL) == "" {
		return nil, fmt.Errorf("short URL is required")
	}

	shortURL := strings.TrimSpace(req.ShortURL)
	if err := validateShortURL(shortURL); err != nil {
		return nil, err
	}

	exists, err := s.repo.ShortURLExists(shortURL, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("short URL %q is already in use", shortURL)
	}

	return s.repo.Update(id, strings.TrimSpace(req.Name), shortURL, strings.TrimSpace(req.TargetURL))
}

func (s *LinkService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *LinkService) HandleRedirect(shortURL, clientIP string) (string, error) {
	link, err := s.repo.GetByShortURL(shortURL)
	if err != nil {
		if err.Error() == fmt.Sprintf("querying link by short_url: %s", sql.ErrNoRows.Error()) {
			return "", fmt.Errorf("link not found")
		}
		return "", err
	}

	if err := s.repo.IncrementClicks(link.ID); err != nil {
		return "", err
	}

	ipHash := hashIP(clientIP)
	isNew, err := s.repo.RecordClick(link.ID, ipHash)
	if err != nil {
		return "", err
	}
	if isNew {
		if err := s.repo.IncrementUniqueClicks(link.ID); err != nil {
			return "", err
		}
	}

	return link.TargetURL, nil
}

func (s *LinkService) generateUniqueShortURL() (string, error) {
	for range maxRetries {
		code, err := generateRandomCode(shortURLLength)
		if err != nil {
			return "", fmt.Errorf("generating random code: %w", err)
		}
		exists, err := s.repo.ShortURLExists(code, 0)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique short URL after %d attempts", maxRetries)
}

func generateRandomCode(length int) (string, error) {
	charsetLen := big.NewInt(int64(len(shortURLCharset)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		b[i] = shortURLCharset[n.Int64()]
	}
	return string(b), nil
}

func validateShortURL(shortURL string) error {
	if len(shortURL) < 2 || len(shortURL) > 64 {
		return fmt.Errorf("short URL must be between 2 and 64 characters")
	}
	for _, c := range shortURL {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return fmt.Errorf("short URL may only contain letters, digits, hyphens, and underscores")
		}
	}
	return nil
}

func hashIP(ip string) string {
	h := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(h[:])
}
