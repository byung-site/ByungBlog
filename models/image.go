package models

type Image struct {
	Model
	Key        string `gorm:"unique_index"`
	Filename   string
	Path       string
	UserID     int
	ArticleKey string
}

func saveImage(image *Image) error {
	return db.Save(image).Error
}
