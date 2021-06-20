package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel `bun:"team,alias:team"`
	DiscordRoleID uint64
	TeamName      string
	CreatedAt     time.Time `bun:"-"`
	ModifiedAt    time.Time `bun:"-"`
}
