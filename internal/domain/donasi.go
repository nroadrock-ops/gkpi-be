package domain
import "time"

type Donasi struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	NamaDonatur   string    `json:"nama_donatur"`
	Jumlah        float64   `json:"jumlah"`
	Status        string    `json:"status"` // pending, success, failed
	PaymentURL    string    `json:"payment_url"` // dari Midtrans
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
