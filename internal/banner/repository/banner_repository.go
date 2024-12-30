package repository

import (
	"log"

	"github.com/akmyrat/global/internal/banner/model"
	"github.com/jmoiron/sqlx"
)

type BannerRepository struct {
	DB *sqlx.DB
}

func NewBannerRepository(db *sqlx.DB) *BannerRepository {
	return &BannerRepository{DB: db}
}

func (r *BannerRepository) Create(banner model.Banner) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction")
		return 0, err
	}

	defer tx.Rollback()

	var bannerID int
	query := `INSERT INTO banner (image_path) VALUES ($1) RETURNING id`
	err = tx.QueryRow(query, banner.ImagePath).Scan(&bannerID)
	if err != nil {
		return 0, err
	}

	for _, translation := range banner.Translations {
		query := `INSERT INTO banner_translate (banner_id, lang_id, title, description) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(query, bannerID, translation.LangID, translation.Title, translation.Description)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (r *BannerRepository) GetAll(langID int) ([]model.BannerByLang, error) {
	var allBanners []model.BannerByLang
	query := `SELECT 
				b.id, l.language, b.image_path, bt.title, bt.description 
			FROM 
				banner AS b 
			INNER JOIN 
				banner_translate AS bt ON b.id=bt.banner_id 
			INNER JOIN 
				languages AS l ON l.id=bt.lang_id 
			WHERE 
				lang_id=$1`

	rows, err := r.DB.Query(query, langID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var banners model.BannerByLang
		err := rows.Scan(
			&banners.BannerID, &banners.LangName, &banners.ImagePath, &banners.Title, &banners.Description,
		)
		if err != nil {
			return nil, err
		}
		allBanners = append(allBanners, banners)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allBanners, nil

}

func (r *BannerRepository) Delete(bannerID int) error {
	query := `DELETE FROM banner WHERE id=$1`
	_, err := r.DB.Exec(query, bannerID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BannerRepository) GetByBannerID(bannerID int) (*model.Banner, error) {
	var banner model.Banner
	translations := []model.Translation{}

	query := `SELECT 
				b.image_path, bt.title, bt.description 
			FROM 
				banner AS b 
			INNER JOIN 
				banner_translate AS bt ON b.id=bt.banner_id 
			INNER JOIN 
				languages AS l ON l.id=bt.lang_id 
			WHERE 
				b.id=$1 `

	rows, err := r.DB.Query(query, bannerID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var translation model.Translation
		err := rows.Scan(
			&banner.ImagePath, &translation.Title, &translation.Description,
		)
		if err != nil {
			return nil, err
		}
		translations = append(translations, translation)
	}
	banner.Translations = translations
	return &banner, nil

}
