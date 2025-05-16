package models

import (
	"gorm.io/gorm"
)

type WaitlistSignup struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Email string `gorm:"not null;unique;uniqueIndex"`
}

func (WaitlistSignup) TableName() string {
	return "waitlist_signups"
}
