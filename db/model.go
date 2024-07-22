package db

import (
	"github.com/google/uuid"
)

type Url struct {
	Id          uuid.UUID  `gorm:"primaryKey" json:"id"`
	OwnerId     *uuid.UUID `gorm:"index" json:"owner_id"`
	Owner       User       `gorm:"foreignKey:OwnerId" json:"-"`
	Path        string     `gorm:"unique" json:"path"`
	RedirectUrl string     `json:"redirect_url"`
}

type User struct {
	Id       uuid.UUID `gorm:"primaryKey" json:"id"`
	Username string    `gorm:"unique" json:"username"`
	Password string    `json:"-"`
}

type UserToken struct {
	Id      uuid.UUID `gorm:"primaryKey" json:"id"`
	UserId  uuid.UUID `gorm:"index" json:"user_id"`
	User    User      `gorm:"foreignKey:UserId" json:"-"`
	Token   string    `gorm:"unique" json:"token"`
	Expires int64     `json:"expires"`
}
