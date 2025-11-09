package model

type Tag struct {
	ID   uint64 `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (Tag) TableName() string {
	return "tags"
}
