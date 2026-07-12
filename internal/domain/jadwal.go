package domain
import "time"

type JadwalIbadah struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	NamaIbadah string   `json:"nama_ibadah"`
	Waktu     time.Time `json:"waktu"`
	Lokasi    string    `json:"lokasi"`
	Pengkhotbah string  `json:"pengkhotbah"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
