package model

import "time"

type Message struct {
	ID          uint `gorm:"primaryKey"`
	SenderID    uint `gorm:"not null"`
	GroupID     *uint
	Content     string `gorm:"not null"`
	ScheduledAt *time.Time
	CreatedAt   time.Time
}

type Group struct {
	ID        uint `gorm:"primaryKey"`
	Name      *string
	CreatedAt time.Time
}

type GroupMember struct {
	ID       uint `gorm:"primaryKey"`
	GroupID  uint
	UserID   uint
	JoinedAt time.Time
}

type BlockedUser struct {
	ID        uint `gorm:"primaryKey"`
	BlockerID uint `gorm:"not null"`
	BlockedID uint `gorm:"not null"`
	CreatedAt time.Time
}
