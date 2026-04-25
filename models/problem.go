package models

import (
	"time"
)

type Problem struct {
	Id                 string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title              string    `gorm:"type:varchar(255);not null" json:"title"`
	Author             string    `gorm:"type:uuid;not null" json:"author"`
	Statement          string    `gorm:"type:text;not null" json:"statement"`
	InputStatement     string    `gorm:"type:text;not null" json:"inputStatement"`
	OutputStatement    string    `gorm:"type:text;not null" json:"outputStatement"`
	TimeLimit          float64   `gorm:"not null" json:"timeLimit"`
	MemoryLimit        float64   `gorm:"not null" json:"memoryLimit"`
	CheckerType        string    `gorm:"type:varchar(10);not null" json:"checkerType"`
	CheckerStrictSpace bool      `gorm:"not null" json:"checkerStrictSpace"`
	CheckerPrecision   *string   `gorm:"type:varchar(10)" json:"checkerPrecision"`
	CreatedAt          time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"type:timestamptz;default:now()" json:"updatedAt"`

	AuthorUser *User `gorm:"foreignKey:Author;references:Id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"authorUser,omitempty"`
}
