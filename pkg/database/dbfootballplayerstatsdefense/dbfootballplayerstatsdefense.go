package dbfootballplayerstatsdefense

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type FootballPlayerStatsDefense struct {
	FootballPlayerStatsDefenseID  int64 `db:"football_player_stats_defense_id"`
	SeasonID                      int64 `db:"season_id"`
	WeekID                        int64 `db:"week_id"`
	PlayerID                      int64 `db:"player_id"`
	OpponentID                    int64 `db:"opponent_id"`
	PointsOwn                     int64 `db:"points_own"`
	PointsOpponent                int64 `db:"points_opponent"`
	TotalYards                    int64 `db:"total_yards"`
	TotalPlays                    int64 `db:"total_plays"`
	PassComp                      int64 `db:"pass_comp"`
	PassAttempts                  int64 `db:"pass_attempts"`
	PassYards                     int64 `db:"pass_yards"`
	PassTDs                       int64 `db:"pass_tds"`
	Sacked                        int64 `db:"sacked"`
	SackedYardsLost               int64 `db:"sacked_yards_lost"`
	PassesDefensed                int64 `db:"passes_defensed"`
	RushAttempts                  int64 `db:"rush_attempts"`
	RushYards                     int64 `db:"rush_yards"`
	RushTDs                       int64 `db:"rush_tds"`
	TacklesForLoss                int64 `db:"tackles_for_loss"`
	TacklesForLossYardsLost       int64 `db:"tackles_for_loss_yards_lost"`
	Interceptions                 int64 `db:"interceptions"`
	InterceptionsReturnYards      int64 `db:"interceptions_return_yards"`
	InterceptionsReturnTDs        int64 `db:"interceptions_return_tds"`
	FumblesForced                 int64 `db:"fumbles_forced"`
	FumbledRecovered              int64 `db:"fumbles_recovered"`
	FumblesReturnYards            int64 `db:"fumbles_return_yards"`
	FumblesReturnTDs              int64 `db:"fumbles_return_tds"`
	Safeties                      int64 `db:"safeties"`
	TwoPtConversions              int64 `db:"two_pt_conversions"`
	TwoPtAttempts                 int64 `db:"two_pt_attempts"`
	Penalties                     int64 `db:"penalties"`
	PenaltiesYards                int64 `db:"penalties_yards"`
	XPMade                        int64 `db:"xp_made"`
	XPAttempts                    int64 `db:"xp_attempts"`
	FGMade                        int64 `db:"fg_made"`
	FGAttempts                    int64 `db:"fg_attempts"`
	FirstDownsTotal               int64 `db:"first_downs_total"`
	FirstDownsRushing             int64 `db:"first_downs_rushing"`
	FirstDownsPassing             int64 `db:"first_downs_passing"`
	FirstDownsPenalty             int64 `db:"first_downs_penalty"`
	ThirdDownConversions          int64 `db:"third_down_conversions"`
	ThirdDownConversionsAttempts  int64 `db:"third_down_conversions_attempts"`
	FourthDownConversions         int64 `db:"fourth_down_conversions"`
	FourthDownConversionsAttempts int64 `db:"fourth_down_conversions_attempts"`
	TimeOfPossession              int64 `db:"time_of_possession"`
	TDsScored                     int64 `db:"tds_scored"`
	PtsScoredByDefense            int64 `db:"pts_scored_by_defense"`
	PtsScoredByOpponentsDefense   int64 `db:"pts_scored_by_opponents_defense"`
	TDDistances                   int64 `db:"td_distances"`
	PtsAllowedFirst               int64 `db:"pts_allowed_first"`
	PtsAllowedSecond              int64 `db:"pts_allowed_second"`
	PtsAllowedThird               int64 `db:"pts_allowed_third"`
	PtsAllowedFourth              int64 `db:"pts_allowed_fourth"`
	PtsAllowedOT                  int64 `db:"pts_allowed_ot"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FootballPlayerStatsDefense, error) {
	d := FootballPlayerStatsDefense{}
	err := database.Get(&d, "SELECT * FROM football_player_stats_defense where division_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FootballPlayerStatsDefense, error) {
	var recs []FootballPlayerStatsDefense
	err := database.Select(&recs, "SELECT * FROM football_player_stats_defense")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FootballPlayerStatsDefense) error {
	_, err := database.Exec("DELETE FROM football_player_stats_defense WHERE football_player_stats_defense_id = ?", d.FootballPlayerStatsDefenseID)
	if err != nil {
		return fmt.Errorf("footballplayerstatsdefense: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FootballPlayerStatsDefense) error {
	res, err := database.Exec(database.BuildInsert("football_player_stats_defense", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("footballplayerstatsdefense: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("footballplayerstatsdefense: couldn't get last inserted ID %S", err)
	}

	d.FootballPlayerStatsDefenseID = ID

	return nil
}
