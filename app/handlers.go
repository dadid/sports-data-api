package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Exception represents an error
type Exception struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func scanPitching(rows *sql.Rows) ([]Pitcher, error) {
	var pitchers []Pitcher
	for rows.Next() {
		var pitcher Pitcher
		err := rows.Scan(
			&pitcher.ID,
			&pitcher.Teamabbrev,
			&pitcher.Rk,
			&pitcher.Pos,
			&pitcher.Name,
			&pitcher.Age,
			&pitcher.W,
			&pitcher.L,
			&pitcher.Wl,
			&pitcher.Era,
			&pitcher.G,
			&pitcher.Gs,
			&pitcher.Gf,
			&pitcher.Cg,
			&pitcher.Sho,
			&pitcher.Sv,
			&pitcher.IP,
			&pitcher.H,
			&pitcher.R,
			&pitcher.Er,
			&pitcher.Hr,
			&pitcher.Bb,
			&pitcher.Ibb,
			&pitcher.So,
			&pitcher.Hbp,
			&pitcher.Bk,
			&pitcher.Wp,
			&pitcher.Bf,
			&pitcher.Eraplus,
			&pitcher.Fip,
			&pitcher.Whip,
			&pitcher.H9,
			&pitcher.Hr9,
			&pitcher.Bb9,
			&pitcher.So9,
			&pitcher.Sow,
			&pitcher.Createddate,
		)
		if err != nil {
			log.Println(errors.Wrap(err, "error scanning row").Error())
		}
		pitchers = append(pitchers, pitcher)
	}
	if err := rows.Err(); err != nil {
		err = errors.Wrap(err, "Error completing query.")
		return nil, err
	}
	return pitchers, nil
}

func scanBatting(rows *sql.Rows) ([]Batter, error) {
	var batters []Batter
	for rows.Next() {
		var batter Batter
		err := rows.Scan(
			&batter.ID,
			&batter.Teamabbrev,
			&batter.Rk,
			&batter.Pos,
			&batter.Name,
			&batter.Age,
			&batter.G,
			&batter.Pa,
			&batter.Ab,
			&batter.R,
			&batter.H,
			&batter.Twob,
			&batter.Threeb,
			&batter.Hr,
			&batter.Rbi,
			&batter.Sb,
			&batter.Cs,
			&batter.Bb,
			&batter.So,
			&batter.Ba,
			&batter.Obp,
			&batter.Slg,
			&batter.Ops,
			&batter.Opsplus,
			&batter.Tb,
			&batter.Gdp,
			&batter.Hbp,
			&batter.Sh,
			&batter.Sf,
			&batter.Ibb,
			&batter.Createddate,
		)
		if err != nil {
			log.Println(errors.Wrap(err, "error scanning row").Error())
		}
		batters = append(batters, batter)
	}
	if err := rows.Err(); err != nil {
		err = errors.Wrap(err, "Error completing query.")
		return nil, err
	}
	return batters, nil
}

func scanBattingSplits(rows *sql.Rows) ([]BattingSplit, error) {
	var battingSplits []BattingSplit
	for rows.Next() {
		var battingSplit BattingSplit
		err := rows.Scan(
			&battingSplit.ID,
			&battingSplit.Teamabbrev,
			&battingSplit.Split,
			&battingSplit.G,
			&battingSplit.Gs,
			&battingSplit.Pa,
			&battingSplit.Ab,
			&battingSplit.R,
			&battingSplit.H,
			&battingSplit.Twob,
			&battingSplit.Threeb,
			&battingSplit.Hr,
			&battingSplit.Rbi,
			&battingSplit.Sb,
			&battingSplit.Cs,
			&battingSplit.Bb,
			&battingSplit.So,
			&battingSplit.Ba,
			&battingSplit.Obp,
			&battingSplit.Slg,
			&battingSplit.Ops,
			&battingSplit.Tb,
			&battingSplit.Gdp,
			&battingSplit.Hbp,
			&battingSplit.Sh,
			&battingSplit.Sf,
			&battingSplit.Ibb,
			&battingSplit.Roe,
			&battingSplit.Babip,
			&battingSplit.Topsplus,
			&battingSplit.Sopsplus,
			&battingSplit.Createddate,
		)
		if err != nil {
			log.Println(errors.Wrap(err, "error scanning row").Error())
		}
		battingSplits = append(battingSplits, battingSplit)
	}
	if err := rows.Err(); err != nil {
		err = errors.Wrap(err, "Error completing query.")
		return nil, err
	}
	return battingSplits, nil
}

func scanPitchingSplits(rows *sql.Rows) ([]PitchingSplit, error) {
	var pitchingSplits []PitchingSplit
	for rows.Next() {
		var pitchingSplit PitchingSplit
		err := rows.Scan(
			&pitchingSplit.ID,
			&pitchingSplit.Teamabbrev,
			&pitchingSplit.Split,
			&pitchingSplit.G,
			&pitchingSplit.Pa,
			&pitchingSplit.Ab,
			&pitchingSplit.R,
			&pitchingSplit.H,
			&pitchingSplit.Twob,
			&pitchingSplit.Threeb,
			&pitchingSplit.Hr,
			&pitchingSplit.Sb,
			&pitchingSplit.Cs,
			&pitchingSplit.Bb,
			&pitchingSplit.So,
			&pitchingSplit.Sow,
			&pitchingSplit.Ba,
			&pitchingSplit.Obp,
			&pitchingSplit.Slg,
			&pitchingSplit.Ops,
			&pitchingSplit.Tb,
			&pitchingSplit.Gdp,
			&pitchingSplit.Hbp,
			&pitchingSplit.Sh,
			&pitchingSplit.Sf,
			&pitchingSplit.Ibb,
			&pitchingSplit.Roe,
			&pitchingSplit.Babip,
			&pitchingSplit.Topsplus,
			&pitchingSplit.Sopsplus,
			&pitchingSplit.Createddate,
		)
		if err != nil {
			log.Println(errors.Wrap(err, "error scanning row").Error())
		}
		pitchingSplits = append(pitchingSplits, pitchingSplit)
	}
	if err := rows.Err(); err != nil {
		err = errors.Wrap(err, "Error completing query.")
		return nil, err
	}
	return pitchingSplits, nil
}

func scanBaserunning(rows *sql.Rows) ([]Baserunner, error) {
	var baserunners []Baserunner
	for rows.Next() {
		var baserunner Baserunner
		err := rows.Scan(
			&baserunner.ID,
			&baserunner.Teamabbrev,
			&baserunner.Name,
			&baserunner.Age,
			&baserunner.Pa,
			&baserunner.Roe,
			&baserunner.Xi,
			&baserunner.Rspct,
			&baserunner.Sbo,
			&baserunner.Sb,
			&baserunner.Cs,
			&baserunner.Sbpct,
			&baserunner.Sb2,
			&baserunner.Cs2,
			&baserunner.Sb3,
			&baserunner.Cs3,
			&baserunner.Sbh,
			&baserunner.Csh,
			&baserunner.Po,
			&baserunner.Pcs,
			&baserunner.Oob,
			&baserunner.Oob1,
			&baserunner.Oob2,
			&baserunner.Oob3,
			&baserunner.Oobhm,
			&baserunner.Bt,
			&baserunner.Xctpct,
			&baserunner.Firsts,
			&baserunner.Firsts2,
			&baserunner.Firsts3,
			&baserunner.Firstd,
			&baserunner.Firstd3,
			&baserunner.Firstdh,
			&baserunner.Seconds,
			&baserunner.Seconds3,
			&baserunner.Secondsh,
			&baserunner.Createddate,
		)
		if err != nil {
			log.Println(errors.Wrap(err, "error scanning row").Error())
		}
		baserunners = append(baserunners, baserunner)
	}
	if err := rows.Err(); err != nil {
		err = errors.Wrap(err, "Error completing query.")
		return nil, err
	}
	return baserunners, nil
}

// GetAllTeams fetches all teams currently in database; endpoint: /api/v1/mlb/teams
func (s *Server) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Query("SELECT id, teamname, teamabbrev FROM baseballreference.team")
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from team table").Error()})
		return
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team
		err := rows.Scan(
			&team.ID,
			&team.TeamName,
			&team.TeamAbbrev,
		)
		if err != nil {
			log.Println(errors.Wrap(err, "error scanning row").Error())
		}
		teams = append(teams, team)
	}
	if err = rows.Err(); err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// GetAllPitching fetches most recent pitching data for all MLB teams; endpoint: /api/v1/mlb/pitching
func (s *Server) GetAllPitching(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, rk, pos, name, age, w, l, wl, era, g, gs, gf, cg, sho, sv, ip, h, r, 
				er, hr, bb, ibb, so, hbp, bk, wp, bf, eraplus, fip, whip, h9, hr9, bb9, so9, sow, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from pitching table").Error()})
		return
	}
	defer rows.Close()

	pitchers, err := scanPitching(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchers)
}

// GetTeamPitching fetches most recent pitching data for specified MLB team; endpoint: /api/v1/mlb/pitching/:teamabbrev
func (s *Server) GetTeamPitching(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, rk, pos, name, age, w, l, wl, era, g, gs, gf, cg, sho, sv, ip, h, r, 
				er, hr, bb, ibb, so, hbp, bk, wp, bf, eraplus, fip, whip, h9, hr9, bb9, so9, sow, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from pitching table").Error()})
		return
	}
	defer rows.Close()

	pitchers, err := scanPitching(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchers)
}

// GetAllBatting fetches most recent batting data for all MLB teams; endpoint: /api/v1/mlb/batting
func (s *Server) GetAllBatting(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, rk, pos, name, age, g, pa, ab, r, h, twob, threeb, hr, rbi, 
				sb, cs, bb, so, ba, obp, slg, ops, opsplus, tb, gdp, hbp, sh, sf, ibb, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting table").Error()})
		return
	}
	defer rows.Close()

	batters, err := scanBatting(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batters)
}

// GetTeamBatting fetches most recent batting data for specified MLB team; endpoint: /api/v1/mlb/batting/:teamabbrev
func (s *Server) GetTeamBatting(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, rk, pos, name, age, g, pa, ab, r, h, twob, threeb, hr, rbi, 
				sb, cs, bb, so, ba, obp, slg, ops, opsplus, tb, gdp, hbp, sh, sf, ibb, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting table").Error()})
		return
	}
	defer rows.Close()

	batters, err := scanBatting(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batters)
}

// GetAllBattingSplits fetches most recent batting_splits data for all MLB teams; endpoint: /api/v1/mlb/splits/batting
func (s *Server) GetAllBattingSplits(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, split, g, gs, pa, ab, r, h, twob, threeb, hr, rbi, sb, cs, bb, so, ba, 
				obp, slg, ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting_splits table").Error()})
		return
	}
	defer rows.Close()

	battingSplits, err := scanBattingSplits(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(battingSplits)
}

// GetTeamBattingSplits fetches most recent batting_splits data for specified MLB team; endpoint: /api/v1/mlb/splits/batting/:teamabbrev
func (s *Server) GetTeamBattingSplits(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, split, g, gs, pa, ab, r, h, twob, threeb, hr, rbi, sb, cs, bb, so, ba, 
				obp, slg, ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting_splits table").Error()})
		return
	}
	defer rows.Close()

	battingSplits, err := scanBattingSplits(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(battingSplits)
}

// GetAllPitchingSplits fetches most recent pitching_splits data for all MLB teams; endpoint: /api/v1/mlb/splits/batting
func (s *Server) GetAllPitchingSplits(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, split, g, pa, ab, r, h, twob, threeb, hr, sb, cs, bb, so, sow, ba, obp, slg, 
				ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting_splits table").Error()})
		return
	}
	defer rows.Close()

	pitchingSplits, err := scanPitchingSplits(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchingSplits)
}

// GetTeamPitchingSplits fetches most recent pitching_splits data for specified MLB team; endpoint: /api/v1/mlb/splits/batting/:teamabbrev
func (s *Server) GetTeamPitchingSplits(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, split, g, pa, ab, r, h, twob, threeb, hr, sb, cs, bb, so, sow, ba, obp, slg, 
				ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting_splits table").Error()})
		return
	}
	defer rows.Close()

	pitchingSplits, err := scanPitchingSplits(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchingSplits)
}

// GetAllBaserunning fetches most recent baserunning data for all MLB teams; endpoint: /api/v1/mlb/splits/baserunning
func (s *Server) GetAllBaserunning(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, name, age, pa, roe, xi, rspct, sbo, sb, cs, sbpct, sb2, cs2, sb3, cs3, sbh, csh, 
				po, pcs, oob, oob1, oob2, oob3, oobhm, bt, xbtpct, firsts, firsts2, firsts3, firstd, firstd3, firstdh, seconds, seconds3, secondsh, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.baserunning p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting_splits table").Error()})
		return
	}
	defer rows.Close()

	baserunning, err := scanBaserunning(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(baserunning)
}

// GetTeamBaserunning fetches most recent baserunning data for specified MLB team; endpoint: /api/v1/mlb/splits/baserunning/:teamabbrev
func (s *Server) GetTeamBaserunning(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	rows, err := s.Dbc.Db.Query(
		`SELECT	id, teamabbrev, name, age, pa, roe, xi, rspct, sbo, sb, cs, sbpct, sb2, cs2, sb3, cs3, sbh, csh, 
				po, pcs, oob, oob1, oob2, oob3, oobhm, bt, xbtpct, firsts, firsts2, firsts3, firstd, firstd3, firstdh, seconds, seconds3, secondsh, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.baserunning p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: errors.Wrap(err, "error selecting from batting_splits table").Error()})
		return
	}
	defer rows.Close()

	baserunning, err := scanBaserunning(rows)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(baserunning)
}
