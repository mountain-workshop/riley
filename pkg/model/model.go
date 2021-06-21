package model

import (
	"time"

	"github.com/uptrace/bun"
)

type AllowedAdmin struct {
	bun.BaseModel  `bun:"allowed_admin"`
	DiscordGuildID uint64 `bun:",notnull"`
	DiscordRoleID  uint64 `bun:",notnull"`
}

type AllowedMod struct {
	bun.BaseModel  `bun:"allowed_mod"`
	DiscordGuildID uint64 `bun:",notnull"`
	DiscordRoleID  uint64 `bun:",notnull"`
}

type League struct {
	bun.BaseModel  `bun:"league"`
	DiscordGuildID uint64 `bun:",notnull"`
	LeagueID       uint32 `bun:",nullzero,notnull,default:'unknown'"`
	FriendlyName   string `bun:",notnull"`
}

type LedgerEntry struct {
	bun.BaseModel     `bun:"ledger_entry"`
	DiscordGuildID    uint64    `bun:",notnull"`
	LedgerEntryID     uint64    `bun:",nullzero,notnull,default:'unknown'"`
	EffectiveDatetime time.Time `bun:",notnull,default:current_timestamp"`
	DiscordUserID     uint64    `bun:",nullzero"`
	TeamID            *Team     `bun:",nullzero,rel:has-one,join:team_id=team_id"`
	TaskID            *Task     `bun:",notnull,rel:has-one,join:task_id=task_id"`
	DiffPointBase     int32     `bun:",notnull"`
	DiffPointBonus    int32     `bun:",notnull"`
}

type ParticipantPointTotal struct {
	bun.BaseModel  `bun:"participant_point_total"`
	DiscordGuildID uint64 `bun:",notnull"`
	DiscordUserID  uint64 `bun:",notnull"`
	PointTotal     int64  `bun:",notnull"`
}

type Season struct {
	bun.BaseModel  `bun:"season"`
	DiscordGuildID uint64 `bun:",notnull"`
	SeasonID       uint32 `bun:",nullzero,notnull,default:'unknown'"`
	FriendlyName   string `bun:",notnull"`
	DatetimeRange  string `bun:",notnull"`
}

type Team struct {
	bun.BaseModel  `bun:"team"`
	DiscordGuildID uint64  `bun:",notnull"`
	TeamID         uint32  `bun:",nullzero,notnull,default:'unknown'"`
	DiscordRoleID  uint64  `bun:",notnull"`
	LeagueID       *League `bun:"rel:has-one,join:league_id=league_id"`
}

type TeamPointTotal struct {
	bun.BaseModel  `bun:"team_point_total"`
	DiscordGuildID uint64 `bun:",notnull"`
	DiscordRoleID  uint64 `bun:",notnull"`
	PointTotal     int64  `bun:",notnull"`
}

type Task struct {
	bun.BaseModel       `bun:"task"`
	DiscordGuildID      uint64                 `bun:",notnull"`
	TaskID              uint64                 `bun:",nullzero,notnull,default:'unknown'"`
	FriendlyName        string                 `bun:",notnull"`
	PointBase           uint32                 `bun:",notnull"`
	PointBonus          uint32                 `bun:",notnull"`
	OpenDatetime        time.Time              `bun:",nullzero,notnull,default:current_timestamp"`
	CloseDatetime       time.Time              `bun:",nullzero"`
	PerParticipantLimit uint16                 `bun:",nullzero"`
	TaskTypeID          *TaskType              `bun:",nullzero"`
	TaskCollectionID    *TaskCollection        `bun:",nullzero"`
	AdditionalInfo      map[string]interface{} `bun:",nullzero"`
}

type TaskCollection struct {
	bun.BaseModel    `bun:"task_collection"`
	DiscordGuildID   uint64 `bun:",notnull"`
	TaskCollectionID uint32 `bun:",nullzero,notnull,default:'unknown'"`
	FriendlyName     string `bun:",notnull"`
}

type TaskParticipantRestriction struct {
	bun.BaseModel                `bun:"task_participant_restriction"`
	DiscordGuildID               uint64    `bun:",notnull"`
	TaskParticipantRestrictionID uint32    `bun:",nullzero,notnull,default:'unknown'"`
	TaskID                       *Task     `bun:",notnull,rel:has-one,join:task_id=task_id"`
	DiscordUserID                uint64    `bun:",notnull"`
	ExpirationDatetime           time.Time `bun:",nullzero"`
}

type TaskType struct {
	bun.BaseModel    `bun:"task_type"`
	DiscordGuildID   uint64    `bun:",notnull"`
	TaskTypeID       uint32    `bun:",nullzero,notnull,default:'unknown'"`
	FriendlyName     string    `bun:",notnull"`
	ParentTaskTypeID *TaskType `bun:",nullzero,rel:has-one,join:task_type_id=task_type_id"`
}
