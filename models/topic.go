package models

type Topic struct {
	Model
	UserId int
	Name   string `gorm:"type:character varying(100)"`
}

//添加或更新话题
func SaveTopic(topic *Topic) error {
	return db.Save(topic).Error
}

//查询所有话题
func QueryTopics() (topics []*Topic, err error) {
	return topics, db.Find(&topics).Error
}

//删除话题
func DeleteTopicById(id int) error {
	return db.Delete(&Topic{}, "id=?", id).Error
}
