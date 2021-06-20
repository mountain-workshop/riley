package model

import (
	"time"

	"github.com/uptrace/bun"
)

type PointsLog struct {
	bun.BaseModel     `bun:"pointslog,alias:pointslog"`
	PointsLogID       uint64 `bun:"primaryKey;<-:false"`
	UserID            uint64
	TeamDiscordRoleID uint64
	Points            uint16
	EffectiveDatetime time.Time `bun:"<-:false"`
	CreatedAt         time.Time `bun:"-"`
	ModifiedAt        time.Time `bun:"-"`
}
