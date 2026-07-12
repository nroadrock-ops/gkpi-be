package domain
import "time"

type Artikel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Judul     string    `json:"judul"`
	Slug      string    `json:"slug" gorm:"unique"`
	Konten    string    `json:"konten"`
	Penulis   string    `json:"penulis"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
