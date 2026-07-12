package domain
import "time"

type Galeri struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Judul     string    `json:"judul"`
	FileURL   string    `json:"file_url"`
	Tags      string    `json:"tags"` // Disimpan dari hasil auto-tagging AI
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
