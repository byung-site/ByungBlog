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

//查询指定用户ID的所有话题
func QueryTopicsByUserID(userId int) (topics []*Topic, err error) {
	return topics, db.Where("user_id=?", userId).Find(&topics).Error
}

//按话题名查询话题
func QueryTopicByName(name string) (topic Topic, err error) {
	return topic, db.Where("name=?", name).Find(&topic).Error
}

//按话题ID查询话题
func QueryTopicByID(id int) (topic Topic, err error) {
	return topic, db.Where("id=?", id).Find(&topic).Error
}

//删除话题
func DeleteTopicById(id int) error {
	return db.Delete(&Topic{}, "id=?", id).Error
}
