package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Exception represents an error
type Exception struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// HandlerFunc overrides htto.HandlerFunc
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func scanData(rows *sqlx.Rows, table string) (data []interface{}, err error) {
	switch table {
	case "pitching":
		for rows.Next() {
			pitcher := Pitcher{}
			err = rows.StructScan(&pitcher)
			if err != nil {
				log.Println(errors.Wrap(err, "error scanning row").Error())
				continue
			}
			data = append(data, pitcher)
		}
	case "batting":
		for rows.Next() {
			batter := Batter{}
			err = rows.StructScan(&batter)
			if err != nil {
				log.Println(errors.Wrap(err, "error scanning row").Error())
				continue
			}
			data = append(data, batter)
		}
	case "baserunning":
		for rows.Next() {
			baserunner := Baserunner{}
			err = rows.StructScan(&baserunner)
			if err != nil {
				log.Println(errors.Wrap(err, "error scanning row").Error())
				continue
			}
			data = append(data, baserunner)
		}
	case "battingsplits":
		for rows.Next() {
			battingSplit := BattingSplit{}
			err = rows.StructScan(&battingSplit)
			if err != nil {
				log.Println(errors.Wrap(err, "error scanning row").Error())
				continue
			}
			data = append(data, battingSplit)
		}
	case "pitchingsplits":
		for rows.Next() {
			pitchingSplit := PitchingSplit{}
			err = rows.StructScan(&pitchingSplit)
			if err != nil {
				log.Println(errors.Wrap(err, "error scanning row").Error())
				continue
			}
			data = append(data, pitchingSplit)
		}
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "Error completing query.")
		return nil, err
	}
	return data, nil
}

// GetAllTeams fetches all teams currently in database; endpoint: /api/v1/mlb/teams
func (s *Server) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams := []Team{}
	err := s.Dbc.Db.Select(&teams, "SELECT id, teamname, teamabbrev FROM baseballreference.team")
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// GetAllPitching fetches most recent pitching data for all MLB teams; endpoint: /api/v1/mlb/pitching
func (s *Server) GetAllPitching(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Queryx(
		`SELECT	id, teamabbrev, rk, pos, name, age, w, l, wl, era, g, gs, gf, cg, sho, sv, ip, h, r, 
				er, hr, bb, ibb, so, hbp, bk, wp, bf, eraplus, fip, whip, h9, hr9, bb9, so9, sow, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	pitchers, err := scanData(rows, "pitching")
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
	rows, err := s.Dbc.Db.Queryx(
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
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	pitchers, err := scanData(rows, "pitching")
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchers)
}

// GetAllBatting fetches most recent batting data for all MLB teams; endpoint: /api/v1/mlb/batting
func (s *Server) GetAllBatting(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Queryx(
		`SELECT	id, teamabbrev, rk, pos, name, age, g, pa, ab, r, h, twob, threeb, hr, rbi, 
				sb, cs, bb, so, ba, obp, slg, ops, opsplus, tb, gdp, hbp, sh, sf, ibb, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	batters, err := scanData(rows, "batting")
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
	rows, err := s.Dbc.Db.Queryx(
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
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	batters, err := scanData(rows, "batting")
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batters)
}

// GetAllBattingSplits fetches most recent batting_splits data for all MLB teams; endpoint: /api/v1/mlb/splits/batting
func (s *Server) GetAllBattingSplits(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Queryx(
		`SELECT	id, teamabbrev, split, g, gs, pa, ab, r, h, twob, threeb, hr, rbi, sb, cs, bb, so, ba, 
				obp, slg, ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	battingSplits, err := scanData(rows, "battingsplits")
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
	rows, err := s.Dbc.Db.Queryx(
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
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	battingSplits, err := scanData(rows, "battingsplits")
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(battingSplits)
}

// GetAllPitchingSplits fetches most recent pitching_splits data for all MLB teams; endpoint: /api/v1/mlb/splits/batting
func (s *Server) GetAllPitchingSplits(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Queryx(
		`SELECT	id, teamabbrev, split, g, pa, ab, r, h, twob, threeb, hr, sb, cs, bb, so, sow, ba, obp, slg, 
				ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	pitchingSplits, err := scanData(rows, "pitchingsplits")
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
	rows, err := s.Dbc.Db.Queryx(
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
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	pitchingSplits, err := scanData(rows, "pitchingsplits")
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchingSplits)
}

// GetAllBaserunning fetches most recent baserunning data for all MLB teams; endpoint: /api/v1/mlb/splits/baserunning
func (s *Server) GetAllBaserunning(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Dbc.Db.Queryx(
		`SELECT	id, teamabbrev, name, age, pa, roe, xi, rspct, sbo, sb, cs, sbpct, sb2, cs2, sb3, cs3, sbh, csh, 
				po, pcs, oob, oob1, oob2, oob3, oobhm, bt, xbtpct, firsts, firsts2, firsts3, firstd, firstd3, firstdh, seconds, seconds3, secondsh, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.baserunning p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	baserunning, err := scanData(rows, "baserunning")
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
	rows, err := s.Dbc.Db.Queryx(
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
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	baserunning, err := scanData(rows, "baserunning")
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(baserunning)
}
