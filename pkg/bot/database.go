package bot

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// DB is the global variable to hold the database connection
// this is done in order to make using the command handlers easier
var DB *gorm.DB

func InitDatabase(host, user, password, dbName, sslmode string, port int) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d, sslmode=%s",
		host,
		user,
		password,
		dbName,
		port,
		sslmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func logPoints(userID, teamID uint64, points uint16) error {
	log := PointsLog{
		UserID:            userID,
		TeamDiscordRoleID: teamID,
		Points:            points,
	}
	ret := DB.Create(&log)
	if ret.Error != nil {
		return fmt.Errorf("error logging points: %v", ret.Error)
	}
	if ret.RowsAffected != 1 {
		return fmt.Errorf("creating points logs affected %d rows", ret.RowsAffected)
	}
	return nil
}

func registerTeam(name string, roleID uint64) error {
	team := Team{
		DiscordRoleID: roleID,
		TeamName:      name,
	}

	ret := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "discord_role_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"team_name"}),
	}).Create(&team)

	if ret.Error != nil {
		return fmt.Errorf("error creating team: %v", ret.Error)
	}
	if ret.RowsAffected != 1 {
		return fmt.Errorf("creating team affected %d rows", ret.RowsAffected)
	}
	return nil
}
