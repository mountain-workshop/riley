package bot

import (
	"time"
)

type Team struct {
	DiscordRoleID uint64 `gorm:"primaryKey"`
	TeamName      string
	CreatedAt     time.Time `gorm:"-"`
	ModifiedAt    time.Time `gorm:"-"`
}

type PointsLog struct {
	PointsLogID       uint64 `gorm:"primaryKey;<-:false"`
	UserID            uint64
	TeamDiscordRoleID uint64
	Points            uint16
	EffectiveDatetime time.Time `gorm:"<-:false"`
	CreatedAt         time.Time `gorm:"-"`
	ModifiedAt        time.Time `gorm:"-"`
}
