package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrSlugEmpty is when generated slug (from title) is emp
var ErrSlugEmpty = errors.New("title must contain at least one alphanumeric character")

var (
	slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)
	slugUUIDPattern  = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type PostContent struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Post is the canonical content entity behind /posts/{slug}.
type Post struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	Slug  string `gorm:"type:text;not null;uniqueIndex"  json:"slug"`
	Title string `gorm:"type:text;not null"  json:"title"`
	Body  string `gorm:"type:text;not null"  json:"body"`

	// Auto-managed by GORM.
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // soft delete - auto-filtered on queries
}

func (p Post) Content() PostContent {
	return PostContent{Slug: p.Slug, Title: p.Title, Body: p.Body}
}

// GenerateSlug produces a URL-safe base slug from a title. Returns empty
// string if the title has no alphanumeric characters (caller should reject).
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = slugInvalidChars.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

// PickAvailableSlug returns the first unused slug by appending -2, -3, ….
func PickAvailableSlug(db *gorm.DB, title string, selfID uuid.UUID) (string, error) {
	base := GenerateSlug(title)

	if base == "" {
		return "", ErrSlugEmpty
	}

	if slugUUIDPattern.MatchString(base) {
		base = "post-" + base
	}

	candidate := base
	for i := 2; ; i++ {
		var count int64
		if err := db.Model(&Post{}).
			Where("slug = ? AND id != ?", candidate, selfID).
			Count(&count).Error; err != nil {
			return "", fmt.Errorf("checking slug availability for %q: %w", candidate, err)
		}
		if count == 0 {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}
