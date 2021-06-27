package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/uptrace/bun"
	"k8s.io/klog"
)

// Team is a roleID that can collect points and have members
type Team struct {
	bun.BaseModel  `bun:"team"`
	DiscordGuildID uint64         `bun:",pk,notnull"`
	TeamID         uint32         `bun:",pk,nullzero,notnull"`
	DiscordRoleID  uint64         `bun:",notnull"`
	LeagueID       uint32         `bun:",nullzero"`
	League         *League        `bun:"rel:belongs-to,join:discord_guild_id=discord_guild_id,league_id=league_id"`
	LedgerEntry    []*LedgerEntry `bun:"rel:has-many,join:discord_guild_id=discord_guild_id,team_id=team_id"`
}

// registerTeamCommand defines the /register-team command
var registerTeamCommand = discordgo.ApplicationCommand{
	Name:        "register-team",
	Description: "Registers a Discord Role as a Team in the tracker",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "discord-role",
			Description: "The role to associate with this team",
			Required:    true,
		},
	},
}

// registerTeamHandler is the slash command handler for /register-team
func registerTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling register-team")

	roleID, err := strconv.ParseUint(i.Data.Options[0].RoleValue(nil, "").ID, 10, 64)
	if err != nil {
		klog.Error(err)
		return
	}

	guildID, err := strconv.ParseUint(i.Interaction.GuildID, 10, 64)
	if err != nil {
		klog.Error(err)
		return
	}

	margs := []interface{}{
		i.Data.Options[0].RoleValue(nil, "").ID,
	}

	var msgFormat string

	_, exists, err := app.createTeam(guildID, roleID)
	if exists {
		msgFormat = " Team <@&%s> is already registered"
	} else if err != nil {
		msgFormat = " Unknown error registering team <@&%s>"
	} else {
		msgFormat = " Team <@&%s> successfully registered\n"
	}

	respond(s, i.Interaction, fmt.Sprintf(msgFormat, margs...))
}

//listTeamCommand defines the /list-teams slash command
var listTeamCommand = discordgo.ApplicationCommand{
	Name:        "list-teams",
	Description: "List all roles associated with teams in the tracker",
}

// listTeamhandler is the slash command handler for /list-teams
func listTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling list-teams")
	var returnMessage string

	guildID, err := strconv.ParseUint(i.Interaction.GuildID, 10, 64)
	if err != nil {
		klog.Error(err)
		returnMessage = " An unexpected error occurred getting your guild id"
		respond(s, i.Interaction, returnMessage)
		return
	}

	teams, err := app.listTeams(guildID)
	if err != nil {
		klog.Error(err)
		returnMessage = " An unexpected error occurred listing teams"
	} else {
		if len(teams) < 1 {
			returnMessage = "There are no teams registered"
		} else {
			returnMessage = "Here are the currently registered teams:\n"
			for _, team := range teams {
				returnMessage = returnMessage + fmt.Sprintf("    <@&%d>\n", team.DiscordRoleID)
			}
		}
	}

	respond(s, i.Interaction, returnMessage)
}

// listTeams lists all the teams for a specific guild
func (app *App) listTeams(guildID uint64) ([]Team, error) {
	teams := make([]Team, 0)
	if err := app.DB().NewSelect().Model(&teams).Where("discord_guild_id = ?", guildID).Scan(app.ctx); err != nil {
		return nil, err
	}
	return teams, nil
}

// getTeam returns a registered team given the guildID and roleID
func (app *App) getTeam(guildID, roleID uint64) (*Team, error) {
	var team Team
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
func (app *App) createTeam(guildID, roleID uint64) (*Team, bool, error) {
	exists := false
	_, err := app.DB().NewInsert().Model(&Team{
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
