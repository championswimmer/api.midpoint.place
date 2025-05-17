package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupPlace represents a suggested place for a group
type GroupPlace struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GroupId   string    `gorm:"not null;uniqueIndex:idx_group_place"`
	PlaceId   string    `gorm:"not null;uniqueIndex:idx_group_place"`
	Name      string    `gorm:"not null"`
	Address   string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Rating    float64   `gorm:"not null"`
	MapURI    string    `gorm:"not null"`
	Latitude  float64   `gorm:"not null"`
	Longitude float64   `gorm:"not null"`
}
