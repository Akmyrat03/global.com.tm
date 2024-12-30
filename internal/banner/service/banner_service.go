package service

import (
	"github.com/akmyrat/global/internal/banner/model"
	"github.com/akmyrat/global/internal/banner/repository"
)

type BannerService struct {
	repo *repository.BannerRepository
}

func NewBannerService(repo *repository.BannerRepository) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) Create(banner model.Banner) (int, error) {
	return s.repo.Create(banner)
}

func (s *BannerService) GetAll(langID int) ([]model.BannerByLang, error) {
	return s.repo.GetAll(langID)
}

func (s *BannerService) Delete(bannerID int) error {
	return s.repo.Delete(bannerID)
}

func (s *BannerService) GetByID(bannerID int) (*model.Banner, error) {
	return s.repo.GetByBannerID(bannerID)
}
