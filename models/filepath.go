package models

type Filepath struct {
	ID        string   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ContestID string   `gorm:"type:uuid;uniqueIndex;not null" json:"contestId"`
	FilePath  string   `gorm:"type:text;not null" json:"filePath"`
	Contest   *Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"contest,omitempty"`
}
