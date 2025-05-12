package models

import "gorm.io/gorm"

// GroupUser represents a many-to-many relationship between users and groups
// It includes the user's location specific to this group membership
type GroupUser struct {
	gorm.Model
	UserID    uint    `gorm:"not null;primaryKey"`
	User      User    `gorm:"foreignKey:UserID"`
	GroupID   string  `gorm:"type:uuid;not null;primaryKey"`
	Group     Group   `gorm:"foreignKey:GroupID"`
	Latitude  float64 `gorm:"type:decimal(10,8);not null"`
	Longitude float64 `gorm:"type:decimal(11,8);not null"`
}

func (GroupUser) TableName() string {
	return "group_users"
}
