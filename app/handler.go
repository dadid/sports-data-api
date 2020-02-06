package app

import (
	"encoding/json"
	"net/http"
)

// Exception represents an error
type Exception struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func checkWriteError(err error, status int, w http.ResponseWriter) bool {
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
		return true
	}
	return false
}

// GetAllTeams fetches all teams currently in database; endpoint: /api/v1/mlb/teams
func (s *Server) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams := []Team{}
	err := s.Dbc.Db.Select(&teams, "SELECT id, teamname, teamabbrev FROM baseballreference.team")
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// GetAllPitching fetches most recent pitching data for all MLB teams; endpoint: /api/v1/mlb/pitching
func (s *Server) GetAllPitching(w http.ResponseWriter, r *http.Request) {
	pitchers := []Pitcher{}
	err := s.Dbc.Db.Select(&pitchers,
		`SELECT	id, teamabbrev, rk, pos, name, age, w, l, wl, era, g, gs, gf, cg, sho, sv, ip, h, r, 
				er, hr, bb, ibb, so, hbp, bk, wp, bf, eraplus, fip, whip, h9, hr9, bb9, so9, sow, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchers)
}

// GetTeamPitching fetches most recent pitching data for specified MLB team; endpoint: /api/v1/mlb/pitching/:teamabbrev
func (s *Server) GetTeamPitching(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	pitchers := []Pitcher{}
	err := s.Dbc.Db.Select(&pitchers,
		`SELECT	id, teamabbrev, rk, pos, name, age, w, l, wl, era, g, gs, gf, cg, sho, sv, ip, h, r, 
				er, hr, bb, ibb, so, hbp, bk, wp, bf, eraplus, fip, whip, h9, hr9, bb9, so9, sow, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchers)
}

// GetAllBatting fetches most recent batting data for all MLB teams; endpoint: /api/v1/mlb/batting
func (s *Server) GetAllBatting(w http.ResponseWriter, r *http.Request) {
	batters := []Batter{}
	err := s.Dbc.Db.Select(&batters,
		`SELECT	id, teamabbrev, rk, pos, name, age, g, pa, ab, r, h, twob, threeb, hr, rbi, 
				sb, cs, bb, so, ba, obp, slg, ops, opsplus, tb, gdp, hbp, sh, sf, ibb, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batters)
}

// GetTeamBatting fetches most recent batting data for specified MLB team; endpoint: /api/v1/mlb/batting/:teamabbrev
func (s *Server) GetTeamBatting(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	batters := []Batter{}
	err := s.Dbc.Db.Select(&batters,
		`SELECT	id, teamabbrev, rk, pos, name, age, g, pa, ab, r, h, twob, threeb, hr, rbi, 
				sb, cs, bb, so, ba, obp, slg, ops, opsplus, tb, gdp, hbp, sh, sf, ibb, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batters)
}

// GetAllBattingSplits fetches most recent batting_splits data for all MLB teams; endpoint: /api/v1/mlb/splits/batting
func (s *Server) GetAllBattingSplits(w http.ResponseWriter, r *http.Request) {
	battingSplits := []BattingSplit{}
	err := s.Dbc.Db.Select(&battingSplits,
		`SELECT	id, teamabbrev, split, g, gs, pa, ab, r, h, twob, threeb, hr, rbi, sb, cs, bb, so, ba, 
				obp, slg, ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(battingSplits)
}

// GetTeamBattingSplits fetches most recent batting_splits data for specified MLB team; endpoint: /api/v1/mlb/splits/batting/:teamabbrev
func (s *Server) GetTeamBattingSplits(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	battingSplits := []BattingSplit{}
	err := s.Dbc.Db.Select(&battingSplits,
		`SELECT	id, teamabbrev, split, g, gs, pa, ab, r, h, twob, threeb, hr, rbi, sb, cs, bb, so, ba, 
				obp, slg, ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.batting_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(battingSplits)
}

// GetAllPitchingSplits fetches most recent pitching_splits data for all MLB teams; endpoint: /api/v1/mlb/splits/batting
func (s *Server) GetAllPitchingSplits(w http.ResponseWriter, r *http.Request) {
	pitchingSplits := []PitchingSplit{}
	err := s.Dbc.Db.Select(&pitchingSplits,
		`SELECT	id, teamabbrev, split, g, pa, ab, r, h, twob, threeb, hr, sb, cs, bb, so, sow, ba, obp, slg, 
				ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchingSplits)
}

// GetTeamPitchingSplits fetches most recent pitching_splits data for specified MLB team; endpoint: /api/v1/mlb/splits/batting/:teamabbrev
func (s *Server) GetTeamPitchingSplits(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	pitchingSplits := []PitchingSplit{}
	err := s.Dbc.Db.Select(&pitchingSplits,
		`SELECT	id, teamabbrev, split, g, pa, ab, r, h, twob, threeb, hr, sb, cs, bb, so, sow, ba, obp, slg, 
				ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.pitching_splits p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pitchingSplits)
}

// GetAllBaserunning fetches most recent baserunning data for all MLB teams; endpoint: /api/v1/mlb/splits/baserunning
func (s *Server) GetAllBaserunning(w http.ResponseWriter, r *http.Request) {
	baserunning := []Baserunner{}
	err := s.Dbc.Db.Select(&baserunning,
		`SELECT	id, teamabbrev, name, age, pa, roe, xi, rspct, sbo, sb, cs, sbpct, sb2, cs2, sb3, cs3, sbh, csh, 
				po, pcs, oob, oob1, oob2, oob3, oobhm, bt, xbtpct, firsts, firsts2, firsts3, firstd, firstd3, firstdh, seconds, seconds3, secondsh, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.baserunning p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
				) x
		WHERE rnk = 1`)
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(baserunning)
}

// GetTeamBaserunning fetches most recent baserunning data for specified MLB team; endpoint: /api/v1/mlb/splits/baserunning/:teamabbrev
func (s *Server) GetTeamBaserunning(w http.ResponseWriter, r *http.Request) {
	p := getParams(r.Context())
	baserunning := []Baserunner{}
	err := s.Dbc.Db.Select(&baserunning,
		`SELECT	id, teamabbrev, name, age, pa, roe, xi, rspct, sbo, sb, cs, sbpct, sb2, cs2, sb3, cs3, sbh, csh, 
				po, pcs, oob, oob1, oob2, oob3, oobhm, bt, xbtpct, firsts, firsts2, firsts3, firstd, firstd3, firstdh, seconds, seconds3, secondsh, createddate
		FROM 	(
					SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
					FROM	baseballreference.baserunning p
							INNER JOIN baseballreference.team t ON t.id = p.teamid
					WHERE	t.teamabbrev = $1
				) x
		WHERE rnk = 1`, p.ByName("teamabbrev"))
	if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(baserunning)
}
