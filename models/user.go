package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string    `json:"username" gorm:"unique;not null;size:50"`
	Password       string    `json:"password" gorm:"not null"`
	RefreshToken   string    `json:"-"`
	LastConnection time.Time `json:"lastConnection"`
	mmm            string
}
