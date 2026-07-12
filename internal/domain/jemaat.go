package domain
import "time"

type Jemaat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	NamaLengkap string  `json:"nama_lengkap"`
	Alamat      string  `json:"alamat"`
	NoTelepon   string  `json:"no_telepon"`
	Status      string  `json:"status"` // aktif, pindah, meninggal
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
