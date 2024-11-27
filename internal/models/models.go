package models

import (
	"gorm.io/gorm"
)

// Config bytes in order;
// ?????

type User struct {
	gorm.Model
	Name           string
	NicknameLocked bool
	SocialCredit   int64
	DiscordID      string `gorm:"unique"`
}

type AudioFile struct {
	gorm.Model
	OwnerID     int
	Owner       User
	Title       string
	Album       *string
	Author      *string
	Length      uint
	FileType    string
	ReleaseYear *uint
	IsPublic    bool
}

type Guild struct {
	gorm.Model
	Name        string
	DjRoles     string
	LoopEnabled bool
}
