package model

import (
	"time"

	"github.com/uptrace/bun"
)

type AllowedAdmin struct {
	bun.BaseModel  `bun:"allowed_admin"`
	DiscordGuildID uint64 `bun:",pk,notnull"`
	DiscordRoleID  uint64 `bun:",pk,notnull"`
}

type AllowedMod struct {
	bun.BaseModel  `bun:"allowed_mod"`
	DiscordGuildID uint64 `bun:",pk,notnull"`
	DiscordRoleID  uint64 `bun:",pk,notnull"`
}

type League struct {
	bun.BaseModel  `bun:"league"`
	DiscordGuildID uint64  `bun:",pk,notnull"`
	LeagueID       uint32  `bun:",pk,nullzero,notnull,default:'unknown'"`
	FriendlyName   string  `bun:",notnull"`
	Team           []*Team `bun:"rel:has-many,join:discord_guild_id=discord_guild_id,league_id=league_id"`
}

type LedgerEntry struct {
	bun.BaseModel     `bun:"ledger_entry"`
	DiscordGuildID    uint64    `bun:",pk,notnull"`
	LedgerEntryID     uint64    `bun:",pk,nullzero,notnull,default:'unknown'"`
	EffectiveDatetime time.Time `bun:",notnull,default:current_timestamp"`
	DiscordUserID     uint64    `bun:",nullzero"`
	TeamID            uint32    `bun:",nullzero"`
	Team              *Team     `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,team_id=team_id"`
	TaskID            uint64    `bun:",notnull"`
	Task              *Task     `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,task_id=task_id"`
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
	DiscordGuildID uint64 `bun:",pk,notnull"`
	SeasonID       uint32 `bun:",pk,nullzero,notnull,default:'unknown'"`
	FriendlyName   string `bun:",notnull"`
	DatetimeRange  string `bun:",notnull"`
}

type Task struct {
	bun.BaseModel              `bun:"task"`
	DiscordGuildID             uint64                        `bun:",pk,notnull"`
	TaskID                     uint64                        `bun:",pk,nullzero,notnull,default:'unknown'"`
	FriendlyName               string                        `bun:",notnull"`
	PointBase                  uint32                        `bun:",notnull"`
	PointBonus                 uint32                        `bun:",notnull"`
	OpenDatetime               time.Time                     `bun:",nullzero,notnull,default:current_timestamp"`
	CloseDatetime              time.Time                     `bun:",nullzero"`
	PerParticipantLimit        uint16                        `bun:",nullzero"`
	TaskTypeID                 uint64                        `bun:",nullzero"`
	TaskType                   *TaskType                     `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,task_type_id=task_type_id"`
	TaskCollectionID           uint64                        `bun:",nullzero"`
	TaskCollection             *TaskCollection               `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,task_collection_id=task_collection_id"`
	AdditionalInfo             map[string]interface{}        `bun:",nullzero"`
	LedgerEntry                []*LedgerEntry                `bun:"rel:has-many,join:discord_guild_id=discord_guild_id,task_id=task_id"`
	TaskParticipantRestriction []*TaskParticipantRestriction `bun:"rel:has-many,join:discord_guild_id=discord_guild_id,task_participant_restriction_id=task_participant_restriction_id"`
}

type TaskCollection struct {
	bun.BaseModel    `bun:"task_collection"`
	DiscordGuildID   uint64  `bun:",pk,notnull"`
	TaskCollectionID uint32  `bun:",pk,nullzero,notnull,default:'unknown'"`
	FriendlyName     string  `bun:",notnull"`
	Task             []*Task `bun:"rel:has-many,join:discord_guild_id=discord_guild_id,task_id=task_id"`
}

type TaskParticipantRestriction struct {
	bun.BaseModel                `bun:"task_participant_restriction"`
	DiscordGuildID               uint64    `bun:",pk,notnull"`
	TaskParticipantRestrictionID uint32    `bun:",pk,nullzero,notnull,default:'unknown'"`
	TaskID                       uint64    `bun:",notnull"`
	Task                         *Task     `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,task_id=task_id"`
	DiscordUserID                uint64    `bun:",notnull"`
	ExpirationDatetime           time.Time `bun:",nullzero"`
}

type TaskType struct {
	bun.BaseModel    `bun:"task_type"`
	DiscordGuildID   uint64    `bun:",pk,notnull"`
	TaskTypeID       uint32    `bun:",pk,nullzero,notnull,default:'unknown'"`
	FriendlyName     string    `bun:",notnull"`
	ParentTaskTypeID uint64    `bun:",nullzero"`
	ParentTaskType   *TaskType `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,parent_task_type_id=task_type_id"`
	Task             []*Task   `bun:"rel:has-many,join:discord_guild_id=discord_guild-id,task_id=task_id"`
}

type Team struct {
	bun.BaseModel  `bun:"team"`
	DiscordGuildID uint64         `bun:",pk,notnull"`
	TeamID         uint32         `bun:",pk,nullzero,notnull,default:'unknown'"`
	DiscordRoleID  uint64         `bun:",notnull"`
	LeagueID       uint32         `bun:",nullzero"`
	League         *League        `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,league_id=league_id"`
	LedgerEntry    []*LedgerEntry `bun:"rel:has-many,join:discord_guild_id=discord_guild_id,team_id=team_id"`
}

type TeamPointTotal struct {
	bun.BaseModel  `bun:"team_point_total"`
	DiscordGuildID uint64 `bun:",notnull"`
	DiscordRoleID  uint64 `bun:",notnull"`
	PointTotal     int64  `bun:",notnull"`
}
