package domain
import "time"

type Absensi struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	JemaatID  uint      `json:"jemaat_id"`
	JadwalID  uint      `json:"jadwal_id"`
	WaktuHadir time.Time `json:"waktu_hadir"`
	Metode    string    `json:"metode"` // manual, ai_face_recognition
}
