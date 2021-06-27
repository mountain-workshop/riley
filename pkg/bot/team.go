package bot

import (
	"strings"

	"github.com/mountain-workshop/riley/pkg/model"
)

func (app *App) listTeams(guildID uint64) ([]model.Team, error) {
	teams := make([]model.Team, 0)
	if err := app.DB().NewSelect().Model(&teams).Where("discord_guild_id = ?", guildID).Scan(app.ctx); err != nil {
		return nil, err
	}
	return teams, nil
}

// getTeam returns a registered team given the guildID and roleID
func (app *App) getTeam(guildID, roleID uint64) (*model.Team, error) {
	var team model.Team
	if err := app.DB().NewSelect().
		Model(&team).
		Where("discord_role_id = ?", roleID).
		Where("discord_guild_id = ?", guildID).
		Scan(app.ctx); err != nil {
		return nil, err
	}
	return &team, nil
}

// createTeam creates a team and returns it. If the team already exists, it returns
// the team and a true boolean indicating that the team already existed.
func (app *App) createTeam(guildID, roleID uint64) (*model.Team, bool, error) {
	exists := false
	_, err := app.DB().NewInsert().Model(&model.Team{
		DiscordRoleID:  roleID,
		DiscordGuildID: guildID,
	}).Exec(app.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			// Team Exists - Get and return it with the boolean set to true
			exists = true
		} else {
			return nil, false, err
		}
	}
	// We created the team. Now get it and return the boolean as false
	team, err := app.getTeam(guildID, roleID)
	return team, exists, err
}
