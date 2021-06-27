package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_listTeams(t *testing.T) {
	_, app := StartTestApp(t)
	fixture := loadFixture(t, app)
	teamOne := fixture.MustRow("Team.one").(*Team)

	tests := []struct {
		name    string
		guildID uint64
		want    []Team
		wantErr bool
	}{
		{
			name:    "one",
			guildID: 1,
			want:    []Team{*teamOne},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := app.listTeams(tt.guildID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestApp_getTeam(t *testing.T) {
	_, app := StartTestApp(t)
	fixture := loadFixture(t, app)
	teamOne := fixture.MustRow("Team.one").(*Team)

	tests := []struct {
		name    string
		guildID uint64
		roleID  uint64
		want    *Team
		wantErr bool
	}{
		{
			name:    "one",
			guildID: 1,
			roleID:  1,
			want:    teamOne,
			wantErr: false,
		},
		{
			name:    "roleID does not exist",
			guildID: 1,
			roleID:  1000,
			wantErr: true,
		},
		{
			name:    "guildID does not exist",
			guildID: 1000,
			roleID:  1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := app.getTeam(tt.guildID, tt.roleID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
