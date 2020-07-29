package app

import (
	"encoding/json"
	"net/http"
	"time"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

// Exception represents an error
type Exception struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func checkWriteError(err error, status int, w http.ResponseWriter) bool {
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Status: status, Message: err.Error()})
		return true
	}
	return false
}

// GenerateToken validates API user creds and returns a JWT token string; endpoint - /user/generateToken
func (s *Server) GenerateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&user)
		if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
			return
		}
		if s.validateCredentials(user) {
			expirationTime := time.Now().Add(72 * time.Hour)
			// Create the JWT claims, which includes the username, password and expiration time
			claims := &Claims{
				Username: user.Username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(), // JWT expiration time is expressed as unix milliseconds
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
			tokenString, err := token.SignedString([]byte(secretKey))
			if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
				return
			}
			// http.SetCookie(w, &http.Cookie{
			// 	Name:    "token",
			// 	Value:   tokenString,
			// 	Expires: expirationTime,
			// })
			json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
			return
		}
		json.NewEncoder(w).Encode(Exception{Status: http.StatusUnauthorized, Message: "error validating credentials"})
	}
}

// GetTeams fetches all teams currently in database; endpoint: /api/v1/mlb/teams
func (s *Server) GetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams := []Team{}
		err := s.Dbc.Db.Select(&teams, "SELECT id, teamname, teamabbrev FROM baseballreference.team")
		if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teams)
	}
}

// GetPitching fetches most recent pitching data for all teams or a specified MLB team; endpoint: /api/v1/mlb/pitching/{teamabbrev}
func (s *Server) GetPitching() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		team := chi.URLParam(r, "teamabbrev")
		pitchers := []Pitcher{}
		var err error
		if team != "" {
			err = s.Dbc.Db.Select(
				&pitchers,
				`SELECT	id, teamabbrev, rk, pos, name, age, w, l, wl, era, g, gs, gf, cg, sho, sv, ip, h, r, 
						er, hr, bb, ibb, so, hbp, bk, wp, bf, eraplus, fip, whip, h9, hr9, bb9, so9, sow, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.pitching p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
							WHERE	t.teamabbrev = $1
						) x
				WHERE rnk = 1`, 
				team,
			)
		} else {
			err = s.Dbc.Db.Select(
				&pitchers,
				`SELECT	id, teamabbrev, rk, pos, name, age, w, l, wl, era, g, gs, gf, cg, sho, sv, ip, h, r, 
						er, hr, bb, ibb, so, hbp, bk, wp, bf, eraplus, fip, whip, h9, hr9, bb9, so9, sow, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.pitching p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
						) x
				WHERE rnk = 1`,
			)
		}
		if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pitchers)
	}
}

// GetBatting fetches most recent batting data for all teams or a specified MLB team; endpoint: /api/v1/mlb/batting/{teamabbrev}
func (s *Server) GetBatting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		team := chi.URLParam(r, "teamabbrev")
		batters := []Batter{}
		var err error
		if team != "" {
			err = s.Dbc.Db.Select(
				&batters,
				`SELECT	id, teamabbrev, rk, pos, name, age, g, pa, ab, r, h, twob, threeb, hr, rbi, 
						sb, cs, bb, so, ba, obp, slg, ops, opsplus, tb, gdp, hbp, sh, sf, ibb, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.batting p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
							WHERE	t.teamabbrev = $1
						) x
				WHERE rnk = 1`,
				team,
			)
		} else {
			err = s.Dbc.Db.Select(
				&batters,
				`SELECT	id, teamabbrev, rk, pos, name, age, g, pa, ab, r, h, twob, threeb, hr, rbi, 
						sb, cs, bb, so, ba, obp, slg, ops, opsplus, tb, gdp, hbp, sh, sf, ibb, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.batting p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
						) x
				WHERE rnk = 1`,
			)
		}
		
		if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(batters)
	}
}

// GetBattingSplits fetches most recent batting_splits data for all teams or specified MLB team; endpoint: /api/v1/mlb/splits/batting/:teamabbrev
func (s *Server) GetBattingSplits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		team := chi.URLParam(r, "teamabbrev")
		battingSplits := []BattingSplit{}
		var err error
		if team != "" {
			err = s.Dbc.Db.Select(
				&battingSplits,
				`SELECT	id, teamabbrev, split, g, gs, pa, ab, r, h, twob, threeb, hr, rbi, sb, cs, bb, so, ba, 
						obp, slg, ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.batting_splits p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
							WHERE	t.teamabbrev = $1
						) x
				WHERE rnk = 1`,
				team,
			)
		} else {
			err = s.Dbc.Db.Select(
				&battingSplits,
				`SELECT	id, teamabbrev, split, g, gs, pa, ab, r, h, twob, threeb, hr, rbi, sb, cs, bb, so, ba, 
						obp, slg, ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.batting_splits p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
						) x
				WHERE rnk = 1`,
			)
		}
		if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(battingSplits)
	}
}

// GetPitchingSplits fetches most recent pitching_splits data for all teams or specified MLB team; endpoint: /api/v1/mlb/splits/batting/{teamabbrev}
func (s *Server) GetPitchingSplits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		team := chi.URLParam(r, "teamabbrev")
		pitchingSplits := []PitchingSplit{}
		var err error
			if team != "" {
				err = s.Dbc.Db.Select(
					&pitchingSplits,
					`SELECT	id, teamabbrev, split, g, pa, ab, r, h, twob, threeb, hr, sb, cs, bb, so, sow, ba, obp, slg, 
							ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
					FROM 	(
								SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
								FROM	baseballreference.pitching_splits p
										INNER JOIN baseballreference.team t ON t.id = p.teamid
								WHERE	t.teamabbrev = $1
							) x
					WHERE rnk = 1`,
					team,
				)
			} else {
				err = s.Dbc.Db.Select(
					&pitchingSplits,
					`SELECT	id, teamabbrev, split, g, pa, ab, r, h, twob, threeb, hr, sb, cs, bb, so, sow, ba, obp, slg, 
							ops, tb, gdp, hbp, sh, sf, ibb, roe, babip, topsplus, sopsplus, createddate
					FROM 	(
								SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
								FROM	baseballreference.pitching_splits p
										INNER JOIN baseballreference.team t ON t.id = p.teamid
							) x
					WHERE rnk = 1`,
				)
			}
		if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pitchingSplits)
	}
}

// GetBaserunning fetches most recent baserunning data for all teamsspecified MLB team; endpoint: /api/v1/mlb/splits/baserunning/{teamabbrev}
func (s *Server) GetBaserunning() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		team := chi.URLParam(r, "teamabbrev")
		baserunning := []Baserunner{}
		var err error
		if team != "" {
			err = s.Dbc.Db.Select(&baserunning,
				`SELECT	id, teamabbrev, name, age, pa, roe, xi, rspct, sbo, sb, cs, sbpct, sb2, cs2, sb3, cs3, sbh, csh, 
						po, pcs, oob, oob1, oob2, oob3, oobhm, bt, xbtpct, firsts, firsts2, firsts3, firstd, firstd3, firstdh, seconds, seconds3, secondsh, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.baserunning p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
							WHERE	t.teamabbrev = $1
						) x
				WHERE rnk = 1`,
				team,
			)
		} else {
			err = s.Dbc.Db.Select(
				&baserunning,
				`SELECT	id, teamabbrev, name, age, pa, roe, xi, rspct, sbo, sb, cs, sbpct, sb2, cs2, sb3, cs3, sbh, csh, 
						po, pcs, oob, oob1, oob2, oob3, oobhm, bt, xbtpct, firsts, firsts2, firsts3, firstd, firstd3, firstdh, seconds, seconds3, secondsh, createddate
				FROM 	(
							SELECT	p.*, t.teamabbrev, DENSE_RANK() OVER(PARTITION BY teamid ORDER BY createddate DESC) AS rnk
							FROM	baseballreference.baserunning p
									INNER JOIN baseballreference.team t ON t.id = p.teamid
						) x
				WHERE rnk = 1`,
			)
		}
		if ok := checkWriteError(err, http.StatusInternalServerError, w); ok {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(baserunning)
	}
}
