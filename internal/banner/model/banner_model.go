package model

type Banner struct {
	ImagePath    string
	Translations []Translation
}

type Translation struct {
	Title       string
	Description string
	LangID      int
}
