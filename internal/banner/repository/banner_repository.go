package repository

import "github.com/jmoiron/sqlx"

type BannerRepository struct {
	DB *sqlx.DB
}
