package main

import "gorm.io/gorm"

type Invoice struct {
	gorm.Model
	FromID        uint `gorm:"not null"`
	ToID          uint `gorm:"not null"`
	Voice         []byte
	Words         []string
	WordsNotAudio bool `gorm:"not null"`
	Seen          bool `gorm:"not null"`
}
