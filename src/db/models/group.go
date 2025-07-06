package models

import (
	"github.com/championswimmer/api.midpoint.place/src/config"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	CreatorID uint   `gorm:"not null"`
	Creator   User   `gorm:"foreignKey:CreatorID"`
	Name      string `gorm:"type:varchar(100);not null"`
	ID        string `gorm:"type:uuid;primary_key;"`
	Code      string `gorm:"type:varchar(10);not null;unique;uniqueIndex"`
	//postgresql
	//Secret    string `gorm:"type:varchar(6);not null;check:secret ~ '^[0-9]{1,6}$'"`
	//sqlite
	Secret string `gorm:"type:varchar(6);not null;"`
	// Type of the group
	// public: visible on main page, searchable by name, anyone can join
	// protected: visible on main page, searchable by name, requires secret to join
	// private: only reachable by group code, requires secret to join
	// postgresql
	// Type              config.GroupType `gorm:"type:enum('public','protected','private');not null;default:'public'"`
	// sqlite
	Type              config.GroupType `gorm:"type:varchar(10);not null;check:type in ('public','protected','private');default:'public'"`
	MidpointLatitude  float64          `gorm:"type:decimal(10,8);not null;default:0"`
	MidpointLongitude float64          `gorm:"type:decimal(11,8);not null;default:0"`
	// Radius in meters
	Radius     int               `gorm:"type:integer;not null;default:2000"`
	PlaceTypes []config.PlaceType `gorm:"type:text;serializer:json"`
	Places     []GroupPlace      `gorm:"foreignKey:GroupID"`
	Members    []GroupUser       `gorm:"foreignKey:GroupID"`
}

func (Group) TableName() string {
	return "groups"
}
