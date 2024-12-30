package model

type Banner struct {
	ImagePath    string        `json:"image_path" db:"image_path" binding:"required"`
	Translations []Translation `json:"translations" binding:"required"`
}

type Translation struct {
	Title       string `json:"title" binding:"required" db:"title"`
	Description string `json:"description" binding:"required" db:"description"`
	LangID      int    `json:"lang_id" binding:"required" db:"lang_id"`
}

type BannerByLang struct {
	BannerID    int
	LangName    string
	ImagePath   string
	Title       string
	Description string
}
