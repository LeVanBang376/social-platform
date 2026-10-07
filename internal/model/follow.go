package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID  uuid.UUID `gorm:"column:follower_id;type:uuid;primaryKey" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"column:following_id;type:uuid;primaryKey" json:"following_id"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`

	// Associations
	Follower  User `gorm:"foreignKey:FollowerID;references:UserID" json:"follower"`
	Following User `gorm:"foreignKey:FollowingID;references:UserID" json:"following"`
}

func (Follow) TableName() string {
	return "follows"
}
