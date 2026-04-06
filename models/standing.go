package models

import (
"time"
)

type ContestStanding struct {
ContestID     uint       `gorm:"primaryKey;autoIncrement:false;index:idx_standings_contest,priority:1" json:"contestId"`
UserID        uint       `gorm:"primaryKey;autoIncrement:false" json:"userId"`
Penalty       int        `gorm:"not null;default:0;index:idx_standings_contest,priority:2,sort:asc" json:"penalty"`
SolvedCount   int        `gorm:"not null;default:0" json:"solvedCount"`
WrongAttempts int        `gorm:"not null;default:0" json:"wrongAttempts"`
LastSolvedAt  *time.Time `gorm:"type:timestamptz" json:"lastSolvedAt"`

Contest Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"contest"`
User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

type ContestSolve struct {
ContestID    uint      `gorm:"primaryKey;autoIncrement:false;index:idx_contest_solves_first_blood,priority:1;index:idx_contest_solves_lookup,priority:1" json:"contestId"`
UserID       uint      `gorm:"primaryKey;autoIncrement:false;index:idx_contest_solves_lookup,priority:2" json:"userId"`
ProblemID    uint      `gorm:"primaryKey;autoIncrement:false;index:idx_contest_solves_first_blood,priority:2;index:idx_contest_solves_lookup,priority:3" json:"problemId"`
SolvedAt     time.Time `gorm:"type:timestamptz;not null;index:idx_contest_solves_lookup,priority:4" json:"solvedAt"`
Penalty      int       `gorm:"not null" json:"penalty"`
AttemptCount int       `gorm:"not null;default:1" json:"attemptCount"`
FirstBlood   bool      `gorm:"not null;default:false;index:idx_contest_solves_first_blood,priority:3" json:"firstBlood"`

Contest Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"contest,omitempty"`
User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
}

type ContestUserProblem struct {
ContestID    uint       `gorm:"primaryKey;autoIncrement:false;index:idx_cup_user_contest,priority:1" json:"contestId"`
UserID       uint       `gorm:"primaryKey;autoIncrement:false;index:idx_cup_user_contest,priority:2" json:"userId"`
ProblemID    uint       `gorm:"primaryKey;autoIncrement:false;index:idx_cup_problem,priority:1;index:idx_cup_problem_id" json:"problemId"`
ProblemIndex int        `gorm:"not null;index:idx_cup_user_contest,priority:3,sort:asc;index:idx_cup_problem,priority:2" json:"problemIndex"`
IsSolved     bool       `gorm:"not null;default:false;index:idx_cup_problem_solved" json:"isSolved"`
SolvedAt     *time.Time `gorm:"type:timestamptz;index:idx_cup_problem_solved,priority:2,sort:asc" json:"solvedAt"`
Penalty      int        `gorm:"not null;default:0" json:"penalty"`
AttemptCount int        `gorm:"not null;default:0" json:"attemptCount"`
FirstBlood   bool       `gorm:"not null;default:false" json:"firstBlood"`

Contest Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"contest,omitempty"`
User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
}

type ContestProblemStat struct {
ContestID      uint `gorm:"primaryKey;autoIncrement:false;index:idx_cps_contest,priority:1" json:"contestId"`
ProblemID      uint `gorm:"primaryKey;autoIncrement:false" json:"problemId"`
ProblemIndex   int  `gorm:"not null;index:idx_cps_contest,priority:2,sort:asc" json:"problemIndex"`
SolvedCount    int  `gorm:"not null;default:0" json:"solvedCount"`
AttemptedUsers int  `gorm:"not null;default:0" json:"attemptedUsers"`

Contest Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"contest,omitempty"`
Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
}
