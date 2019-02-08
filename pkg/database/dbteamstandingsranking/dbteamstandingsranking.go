package dbteamstandingsranking

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type TeamStandingsRanking struct {
	TeamStandingsRankingID int64               `db:"team_standings_ranking_id"`
	TeamStandingsID        int64               `db:"team_standings_id"`
	RankingSystemID        int64               `db:"ranking_system_id"`
	Ranking                database.NullString `db:"ranking"`
	DisplayOrder           database.NullInt64  `db:"display_order"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*TeamStandingsRanking, error) {
	d := TeamStandingsRanking{}
	err := database.Get(&d, "SELECT * FROM team_standings_ranking where team_standings_ranking_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]TeamStandingsRanking, error) {
	var recs []TeamStandingsRanking
	err := database.Select(&recs, "SELECT * FROM team_standings_ranking")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *TeamStandingsRanking) error {
	_, err := database.Exec("DELETE FROM team_standings_ranking WHERE team_standings_ranking_id = ?", d.TeamStandingsRankingID)
	if err != nil {
		return fmt.Errorf("teamstandingsranking: couldn't delete team_standings_ranking %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *TeamStandingsRanking) error {
	res, err := database.Exec(database.BuildInsert("team_standings_ranking", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("teamstandingsranking: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("teamstandingsranking: couldn't get last inserted ID %S", err)
	}

	d.TeamStandingsRankingID = ID

	return nil
}
