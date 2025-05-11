package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	Username  string  `gorm:"not null;unique;uniqueIndex"`
	Password  string  `gorm:"not null"`
	Latitude  float64 `gorm:"type:decimal(10,8);not null;default:0"`
	Longitude float64 `gorm:"type:decimal(11,8);not null;default:0"`
}
