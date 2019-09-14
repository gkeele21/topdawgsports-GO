package dbfootballplayerstats

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FootballPlayerStats struct {
	FootballPlayerStatsID int64               `db:"football_player_stats_id"`
	SeasonID              int64               `db:"season_id"`
	WeekID                database.NullInt64  `db:"week_id"`
	PlayerID              int64               `db:"player_id"`
	Started               database.NullInt64  `db:"started"`
	Played                database.NullInt64  `db:"played"`
	Inactive              database.NullInt64  `db:"inactive"`
	PassComp              database.NullInt64  `db:"pass_comp"`
	PassAttempts          database.NullInt64  `db:"pass_attempts"`
	PassYards             database.NullInt64  `db:"pass_yards"`
	PassInterceptions     database.NullInt64  `db:"pass_interceptions"`
	PassTDs               database.NullInt64  `db:"pass_tds"`
	PassTwoPts            database.NullInt64  `db:"pass_twopts"`
	Sacked                database.NullInt64  `db:"sacked"`
	SackedYardsLost       database.NullInt64  `db:"sacked_yards_lost"`
	RushAttempts          database.NullInt64  `db:"rush_attempts"`
	RushYards             database.NullInt64  `db:"rush_yards"`
	RushTDs               database.NullInt64  `db:"rush_tds"`
	RushTwoPts            database.NullInt64  `db:"rush_twopts"`
	RecTargets            database.NullInt64  `db:"rec_targets"`
	RecCatches            database.NullInt64  `db:"rec_catches"`
	RecYards              database.NullInt64  `db:"rec_yards"`
	RecTDs                database.NullInt64  `db:"rec_tds"`
	RecTwoPts             database.NullInt64  `db:"rec_twopts"`
	XPMade                database.NullInt64  `db:"xp_made"`
	XPAttempts            database.NullInt64  `db:"xp_attempts"`
	XPBlocked             database.NullInt64  `db:"xp_blocked"`
	FGMade                database.NullInt64  `db:"fg_made"`
	FGAttempts            database.NullInt64  `db:"fg_attempts"`
	FGBlocked             database.NullInt64  `db:"fg_blocked"`
	FG29Minus             database.NullInt64  `db:"fg_29_minus"`
	FG3039                database.NullInt64  `db:"fg_30_39"`
	FG4049                database.NullInt64  `db:"fg_40_49"`
	FG50Plus              database.NullInt64  `db:"fg_50_plus"`
	Fumbles               database.NullInt64  `db:"fumbles"`
	FumblesLost           database.NullInt64  `db:"fumbles_lost"`
	ExtraTDs              database.NullInt64  `db:"extra_tds"`
	TDDistances           database.NullString `db:"td_distances"`
	DefAssists            database.NullInt64  `db:"def_assists"`
	DefFumbleRecoveries   database.NullInt64  `db:"def_fumble_recoveries"`
	DefFumbleReturnYards  database.NullInt64  `db:"def_fumble_return_yards"`
	DefFumbleReturnForTDs database.NullInt64  `db:"def_fumble_return_for_tds"`
	DefFumbleForced       database.NullInt64  `db:"def_fumble_forced"`
	DefIntReturnYards     database.NullInt64  `db:"def_int_return_yards"`
	DefIntReturnForTDs    database.NullInt64  `db:"def_int_return_for_tds"`
	DefInt                database.NullInt64  `db:"def_int"`
	DefPassesDefensed     database.NullInt64  `db:"def_passes_defensed"`
	DefQBHurries          database.NullInt64  `db:"def_qb_hurries"`
	DefSackYardsLost      database.NullInt64  `db:"def_sack_yards_lost"`
	DefSacks              database.NullInt64  `db:"def_sacks"`
	DefSafeties           database.NullInt64  `db:"def_safeties"`
	DefTDDistances        database.NullString `db:"def_td_distances"`
	DefTackles            database.NullInt64  `db:"def_tackles"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FootballPlayerStats, error) {
	d := FootballPlayerStats{}
	err := database.Get(&d, "SELECT * FROM football_player_stats where football_player_stats_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FootballPlayerStats, error) {
	var recs []FootballPlayerStats
	err := database.Select(&recs, "SELECT * FROM football_player_stats")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FootballPlayerStats) error {
	_, err := database.Exec("DELETE FROM football_player_stats WHERE football_player_stats_id = ?", d.FootballPlayerStatsID)
	if err != nil {
		return fmt.Errorf("footballplayerstats: couldn't delete footballplayerstats %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FootballPlayerStats) error {
	fmt.Printf("Ready to insert record : %#v \n", d)
	res, err := database.Exec(database.BuildInsert("football_player_stats", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("footballplayerstats: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("footballplayerstats: couldn't get last inserted ID %S", err)
	}

	d.FootballPlayerStatsID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FootballPlayerStats) error {
	sql := database.BuildUpdate("football_player_stats", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("football_player_stats: couldn't update %s", err)
	}

	return nil
}

func Save(s *FootballPlayerStats) error {
	if s.FootballPlayerStatsID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadBySeasonWeekPlayer reads by season_id, week_id, and player_id columns
func ReadBySeasonWeekPlayer(seasonId, weekId, playerId int64) (*FootballPlayerStats, error) {
	d := FootballPlayerStats{}
	err := database.Get(&d, "SELECT * FROM football_player_stats where season_id = ? AND week_id = ? AND player_id = ?", seasonId, weekId, playerId)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
