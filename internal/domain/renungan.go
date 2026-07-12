package domain
import "time"

type Renungan struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Judul     string    `json:"judul"`
	Ayat      string    `json:"ayat"`
	Konten    string    `json:"konten"`
	Tanggal   time.Time `json:"tanggal" gorm:"type:date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
