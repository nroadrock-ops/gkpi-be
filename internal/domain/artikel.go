package domain
import "time"

type Artikel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Judul       string    `json:"judul"`
	Slug        string    `json:"slug" gorm:"unique"`
	Konten      string    `json:"konten"`
	Ringkasan   string    `json:"ringkasan"`
	Kategori    string    `json:"kategori" gorm:"default:'berita'"`
	Penulis     string    `json:"penulis"`
	GambarUrl   string    `json:"gambar_url"`
	Dipublikasi bool      `json:"dipublikasi" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
