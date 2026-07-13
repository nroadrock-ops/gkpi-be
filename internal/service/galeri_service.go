package service

import (
	"fmt"
	"log"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type GaleriService interface {
	GetAll() ([]domain.Galeri, error)
	UploadAndSave(judul, fileURL string) (*domain.Galeri, error)
	Delete(id uint) error
}

type galeriService struct {
	repo     repository.GaleriRepository
	aiClient AIServiceClient
}

// Inject repo dan aiClient
func NewGaleriService(repo repository.GaleriRepository, aiClient AIServiceClient) GaleriService {
	return &galeriService{repo, aiClient}
}

func (s *galeriService) GetAll() ([]domain.Galeri, error) {
	return s.repo.FindAll()
}

func (s *galeriService) UploadAndSave(judul, fileURL string) (*domain.Galeri, error) {
	galeri := &domain.Galeri{
		Judul:   judul,
		FileURL: fileURL,
		Tags:    "", // Kosong dulu, diisi AI nanti
	}

	// Simpan ke DB segera agar response tidak lambat
	if err := s.repo.Create(galeri); err != nil {
		return nil, err
	}

	// Panggil AI service secara async (goroutine) untuk auto-tagging
	go func(gID uint, fURL string) {
		log.Printf("Starting async AI auto-tagging for Galeri ID %d", gID)
		
		payload := map[string]string{
			"image_url": fURL,
		}
		
		res, err := s.aiClient.ForwardRequest("/galeri/classify", payload)
		if err != nil {
			// Tangani kasus AI service down dengan graceful error
			log.Printf("[Galeri Async] AI auto-tagging failed for ID %d: %v. Foto tetap tersimpan tanpa tag.", gID, err)
			return
		}

		// Anggap response dari AI ada field "tags" berupa string atau list of string
		var tagsStr string
		if data, ok := res["data"].(map[string]interface{}); ok {
			if tags, ok := data["tags"].(string); ok {
				tagsStr = tags
			} else if tagsArr, ok := data["tags"].([]interface{}); ok {
				for i, t := range tagsArr {
					if i > 0 {
						tagsStr += ", "
					}
					tagsStr += fmt.Sprintf("%v", t)
				}
			}
		}

		if tagsStr != "" {
			// Update tag ke DB
			g, err := s.repo.FindByID(gID)
			if err == nil {
				g.Tags = tagsStr
				_ = s.repo.Update(g)
				log.Printf("[Galeri Async] Successfully updated tags for ID %d: %s", gID, tagsStr)
			}
		}
	}(galeri.ID, galeri.FileURL)

	return galeri, nil
}

func (s *galeriService) Delete(id uint) error {
	// In a real scenario, also delete file from Supabase Storage here
	return s.repo.Delete(id)
}
