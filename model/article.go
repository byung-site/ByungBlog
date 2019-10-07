package model

type Article struct {
	Model
	Key     string `gorm:"unique_index"`
	UserID  int
	User    User  `gorm:-`
	Topic   Topic `gorm:-`
	TopicID int
	Image   string
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
	return article, db.Where("key=?", key).Preload("User").Preload("Topic").Take(&article).Error
}

//查询指定topicid的文章
func QueryArticlesByTopicID(topicid uint) (articles []*Article, err error) {
	return articles, db.Where("topic_id=?", topicid).Order("created_at desc").Preload("User").Preload("Topic").Find(&articles).Error
}

//查询指定userid的文章
func QueryArticlesByUserID(userid uint) (articles []*Article, err error) {
	return articles, db.Where("user_id=?", userid).Order("created_at desc").Preload("User").Preload("Topic").Find(&articles).Error
}

//查询指定topicid的文章数
func QueryArticleCountByTopicID(topicid uint) (count int, err error) {
	return count, db.Table("articles").Where("topic_id=? and deleted_at IS NULL", topicid).Count(&count).Error
}

//查询所有文章
func QueryAllArticles() (articles []*Article, err error) {
	return articles, db.Order("created_at desc").Preload("User").Preload("Topic").Where("deleted_at IS NULL").Find(&articles).Error
}

//查询所有发布的文章
func QueryPublishArticles(userId int) (articles []*Article, err error) {
	return articles, db.Where("publish = ? and user_id = ?", 1, userId).Order("created_at").Preload("User").Preload("Topic").Find(&articles).Error
}

//查询最热的10篇文章
func QueryHottestArticle() (articles []*Article, err error) {
	return articles, db.Limit(10).Where("publish=?", 1).Order("visit desc").Preload("User").Preload("Topic").Find(&articles).Error
}

//查询最新的10篇文章
func QueryNewestArticle() (articles []*Article, err error) {
	return articles, db.Limit(10).Where("publish=?", 1).Order("created_at desc").Preload("User").Preload("Topic").Find(&articles).Error
}

//删除文章
func DeleteArticleByKey(key string) error {
	return db.Delete(&Article{}, "key=?", key).Error
}
