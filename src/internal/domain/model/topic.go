package model

type Topic struct {
	ID   uint64 `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (Topic) TableName() string {
	return "topics"
}
