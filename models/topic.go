package models

type Topic struct {
	Model
	UserId int
	Name   string `gorm:"type:character varying(100)" json:"topicName"`
}

func SaveTopic(topic *Topic) error {
	return db.Save(topic).Error
}

func QueryAllTopics() (topics []*Topic, err error) {
	return topics, db.Find(&topics).Error
}

func DeleteTopicById(id int) error {
	return db.Delete(&Topic{}, "id=?", id).Error
}
