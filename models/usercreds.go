package models

type UserCreds struct {
	Id            string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ContestId     string `gorm:"type:uuid;not null;uniqueIndex:uq_user_creds_contest_user,priority:1" json:"contestId"`
	UserId        string `gorm:"type:uuid;not null;uniqueIndex:uq_user_creds_contest_user,priority:2" json:"userId"`
	PlainPassword string `gorm:"type:varchar(255);not null" json:"password"`

	User    User    `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Contest Contest `gorm:"foreignKey:ContestId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"contest"`
}
