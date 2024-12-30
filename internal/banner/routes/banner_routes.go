package routes

import (
	"github.com/akmyrat/global/internal/banner/middleware"
	"github.com/akmyrat/global/internal/banner/repository"
	"github.com/akmyrat/global/internal/banner/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func BannerRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	bannerRepository := repository.NewBannerRepository(db)
	bannerService := service.NewBannerService(bannerRepository)
	bannerMiddleware := middleware.NewBannerMiddleware(bannerService)

	bannerRoutes := router.Group("/banners")
	{
		bannerRoutes.POST("/", bannerMiddleware.CreateBanner())
		bannerRoutes.DELETE("/:id", bannerMiddleware.DeleteBanner())
	}
	router.GET("/banners", bannerMiddleware.GetAllBannersByLang())
}
