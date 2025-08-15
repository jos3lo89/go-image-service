package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	Filename    string `gorm:"not null"`
	PersonName  string `gorm:"index"`
	PersonDNI   string `gorm:"index"`
	URL         string
	ContentType string
}
