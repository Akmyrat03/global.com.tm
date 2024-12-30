package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/akmyrat/global/internal/banner/model"
	"github.com/akmyrat/global/internal/banner/service"
	handler "github.com/akmyrat/global/pkg/response"
	"github.com/gin-gonic/gin"
)

type BannerMiddleware struct {
	service *service.BannerService
}

func NewBannerMiddleware(service *service.BannerService) *BannerMiddleware {
	return &BannerMiddleware{service: service}
}

func (m *BannerMiddleware) CreateBanner() gin.HandlerFunc {
	return func(c *gin.Context) {
		var banner model.Banner

		image, err := c.FormFile("image_path")
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		uploadDir := "./uploads/banners"

		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, 0755)
		}

		filepath := filepath.Join(uploadDir, image.Filename)

		if err := c.SaveUploadedFile(image, filepath); err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		titleTurkmen := c.PostForm("title_turkmen")
		titleEnglish := c.PostForm("title_english")
		titleRussian := c.PostForm("title_russian")

		descriptionTurkmen := c.PostForm("description_turkmen")
		descriptionEnglish := c.PostForm("description_english")
		descriptionRussian := c.PostForm("description_russian")

		translations := []model.Translation{
			{LangID: 1, Title: titleTurkmen, Description: descriptionTurkmen},
			{LangID: 2, Title: titleEnglish, Description: descriptionEnglish},
			{LangID: 3, Title: titleRussian, Description: descriptionRussian},
		}

		banner = model.Banner{
			ImagePath:    filepath,
			Translations: translations,
		}

		id, err := m.service.Create(banner)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully created banner",
			"id":      id,
		})
	}
}

func (m *BannerMiddleware) GetAllBannersByLang() gin.HandlerFunc {
	return func(c *gin.Context) {
		langID, err := strconv.Atoi(c.Query("lang_id"))
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid language ID")
			return
		}

		banners, err := m.service.GetAll(langID)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"banners": banners,
		})
	}
}

func (m *BannerMiddleware) DeleteBanner() gin.HandlerFunc {
	return func(c *gin.Context) {
		bannerID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		existingBanner, err := m.service.GetByID(bannerID)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusNotFound, "Could not find banner")
			return
		}

		err = m.service.Delete(bannerID)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, "Could not delete banner")
			return
		}

		if existingBanner.ImagePath != "" {
			if err := os.Remove(existingBanner.ImagePath); err != nil && !os.IsNotExist(err) {
				handler.NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete banner image")
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Successfully deleted banner",
			"banner_id": bannerID,
		})

	}
}
