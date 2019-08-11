package models

type Article struct {
	Model
	Key     string `gorm:"unique_index"`
	UserID  int
	TopicID int
	User    User
	Topic   Topic
	Title   string `gorm:"type:character varying(200)"`
	Summary string `gorm:"type:character varying(800)"`
	Content string `gorm:"type:text"`
	Visit   int    `gorm:"default:0"`
	Praise  int    `gorm:"default:0"`
	Publish int    `gorm:"default:0"`
}

func SaveArticle(article *Article) error {
	return db.Save(article).Error
}

func QueryArticleByKey(key string) (article Article, err error) {
	return article, db.Where("key=?", key).Take(&article).Error
}
