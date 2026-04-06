package models

type Filepath struct {
	ID        uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ContestID uint     `gorm:"uniqueIndex;not null" json:"contestId"`
	FilePath  string   `gorm:"type:text;not null" json:"filePath"`
	Contest   *Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"contest,omitempty"`
}
