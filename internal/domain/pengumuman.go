package domain
import "time"

type Pengumuman struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Judul     string    `json:"judul"`
	Konten    string    `json:"konten"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
