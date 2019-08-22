package models

type Article struct {
	Model
	Key     string `gorm:"unique_index"`
	UserID  int
	User    User
	TopicID int
	Title   string `gorm:"type:character varying(200)"`
	Summary string `gorm:"type:character varying(800)"`
	Content string `gorm:"type:text"`
	Visit   int    `gorm:"default:0"`
	Praise  int    `gorm:"default:0"`
	Publish int    `gorm:"default:0"`
}

//保存文章
func SaveArticle(article *Article) error {
	return db.Save(article).Error
}

//查询指定key的文章
func QueryArticleByKey(key string) (article Article, err error) {
	return article, db.Where("key=?", key).Take(&article).Error
}

//查询指定topicid的文章
func QueryArticlesByTopicID(topicid uint) (articles []*Article, err error) {
	return articles, db.Where("topic_id=?", topicid).Find(&articles).Error
}

//查询指定topicid的文章数
func QueryArticleCountByTopicID(topicid uint) (count int, err error) {
	return count, db.Table("articles").Where("topic_id=?", topicid).Count(&count).Error
}

//查询所有文章
func QueryAllArticles() (articles []*Article, err error) {
	return articles, db.Find(&articles).Order("create_at").Error
}

//查询最热的10篇文章
func QueryHottestArticle() (articles []*Article, err error) {
	return articles, db.Limit(10).Find(&articles).Order("visit").Error
}

//查询最新的10篇文章
func QueryNewestArticle() (articles []*Article, err error) {
	return articles, db.Limit(10).Find(&articles).Order("create_at").Error
}

//删除文章
func DeleteArticleByKey(key string) error {
	return db.Delete(&Article{}, "key=?", key).Error
}
