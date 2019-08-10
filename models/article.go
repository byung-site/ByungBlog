package models

type Article struct {
	Model
	Key     string `gorm:"unique:not null"`
	UserID  int
	TopicID int
	User    User
	Topic   Topic
	Title   string `gorm:"type:character varying(200)"`
	Summary string `gorm:"type:character varying(800)"`
	Content string `gorm:"type:text"`
	Visit   int    `gorm:"default:0"`
	Praise  int    `gorm:"default:0"`
}

func SaveArticle(article *Article) error {
	return db.Save(article).Error
}
