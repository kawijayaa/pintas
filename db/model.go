package db

import (
	"github.com/google/uuid"
)

type ShortUrl struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	Path        string    `gorm:"unique" json:"path"`
	RedirectUrl string    `json:"redirect_url"`
}
