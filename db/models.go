package db

import "gorm.io/gorm"

type Invoice struct {
	gorm.Model
	FromID        uint32 `gorm:"not null"`
	ToID          uint32 `gorm:"not null"`
	Voice         []byte
	Words         []string
	WordsNotAudio bool `gorm:"not null"`
	Seen          bool `gorm:"not null"`
}
