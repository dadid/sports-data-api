package app

import "gopkg.in/guregu/null.v3"

// Team represents a single MLB team
type Team struct {
	ID         string `json:"id"`
	TeamName   string `json:"teamname"`
	TeamAbbrev string `json:"teamabbrev"`
}

// Pitcher represents data for a single pitcher
type Pitcher struct {
	ID          int         `json:"id"`
	Teamabbrev  string      `json:"teamabbrev"`
	Rk          int         `json:"rk"`
	Pos         null.String `json:"pos"`
	Name        null.String `json:"name"`
	Age         null.Int    `json:"age"`
	W           null.Int    `json:"w"`
	L           null.Int    `json:"l"`
	Wl          null.Float  `json:"wl"`
	Era         null.Float  `json:"era"`
	G           null.Int    `json:"g"`
	Gs          null.Int    `json:"gs"`
	Gf          null.Int    `json:"gf"`
	Cg          null.Int    `json:"cg"`
	Sho         null.Int    `json:"sho"`
	Sv          null.Int    `json:"sv"`
	IP          null.Float  `json:"ip"`
	H           null.Int    `json:"h"`
	R           null.Int    `json:"r"`
	Er          null.Int    `json:"er"`
	Hr          null.Int    `json:"hr"`
	Bb          null.Int    `json:"bb"`
	Ibb         null.Int    `json:"ibb"`
	So          null.Int    `json:"so"`
	Hbp         null.Int    `json:"hbp"`
	Bk          null.Int    `json:"bk"`
	Wp          null.Int    `json:"wp"`
	Bf          null.Int    `json:"bf"`
	Eraplus     null.Int    `json:"eraplus"`
	Fip         null.Float  `json:"fip"`
	Whip        null.Float  `json:"whip"`
	H9          null.Float  `json:"h9"`
	Hr9         null.Float  `json:"hr9"`
	Bb9         null.Float  `json:"bb9"`
	So9         null.Float  `json:"so9"`
	Sow         null.Float  `json:"sow"`
	Createddate null.Time   `json:"createddate"`
}

// Batter represents data for a single batter
type Batter struct {
	ID          int         `json:"id"`
	Teamabbrev  string      `json:"teamabbrev"`
	Rk          int         `json:"rk"`
	Pos         null.String `json:"pos"`
	Name        null.String `json:"name"`
	Age         null.Int    `json:"age"`
	G           null.Int    `json:"g"`
	Pa          null.Int    `json:"pa"`
	Ab          null.Int    `json:"ab"`
	R           null.Int    `json:"r"`
	H           null.Int    `json:"h"`
	Twob        null.Int    `json:"twob"`
	Threeb      null.Int    `json:"threeb"`
	Hr          null.Int    `json:"hr"`
	Rbi         null.Int    `json:"rbi"`
	Sb          null.Int    `json:"sb"`
	Cs          null.Int    `json:"cs"`
	Bb          null.Int    `json:"bb"`
	So          null.Int    `json:"so"`
	Ba          null.Float  `json:"ba"`
	Obp         null.Float  `json:"obp"`
	Slg         null.Float  `json:"slg"`
	Ops         null.Float  `json:"ops"`
	Opsplus     null.Int    `json:"opsplus"`
	Tb          null.Int    `json:"tb"`
	Gdp         null.Int    `json:"gdp"`
	Hbp         null.Int    `json:"hbp"`
	Sh          null.Int    `json:"sh"`
	Sf          null.Int    `json:"sf"`
	Ibb         null.Int    `json:"ibb"`
	Createddate null.Time   `json:"createddate"`
}

// BattingSplit represents data for a batting_splits
type BattingSplit struct {
	ID          int         `json:"id"`
	Teamabbrev  string      `json:"teamabbrev"`
	Split       null.String `json:"split"`
	G           null.Int    `json:"g"`
	Gs          null.Int    `json:"gs"`
	Pa          null.Int    `json:"pa"`
	Ab          null.Int    `json:"ab"`
	R           null.Int    `json:"r"`
	H           null.Int    `json:"h"`
	Twob        null.Int    `json:"twob"`
	Threeb      null.Int    `json:"threeb"`
	Hr          null.Int    `json:"hr"`
	Rbi         null.Int    `json:"rbi"`
	Sb          null.Int    `json:"sb"`
	Cs          null.Int    `json:"cs"`
	Bb          null.Int    `json:"bb"`
	So          null.Int    `json:"so"`
	Ba          null.Float  `json:"ba"`
	Obp         null.Float  `json:"obp"`
	Slg         null.Float  `json:"slg"`
	Ops         null.Float  `json:"ops"`
	Tb          null.Int    `json:"tb"`
	Gdp         null.Int    `json:"gdp"`
	Hbp         null.Int    `json:"hbp"`
	Sh          null.Int    `json:"sh"`
	Sf          null.Int    `json:"sf"`
	Ibb         null.Int    `json:"ibb"`
	Roe         null.Int    `json:"roe"`
	Babip       null.Float  `json:"babip"`
	Topsplus    null.Int    `json:"topsplus"`
	Sopsplus    null.Int    `json:"sopsplus"`
	Createddate null.Time   `json:"createddate"`
}

// PitchingSplit represents data for a pitching splits
type PitchingSplit struct {
	ID          int         `json:"id"`
	Teamabbrev  string      `json:"teamabbrev"`
	Split       null.String `json:"split"`
	G           null.Int    `json:"g"`
	Pa          null.Int    `json:"pa"`
	Ab          null.Int    `json:"ab"`
	R           null.Int    `json:"r"`
	H           null.Int    `json:"h"`
	Twob        null.Int    `json:"twob"`
	Threeb      null.Int    `json:"threeb"`
	Hr          null.Int    `json:"hr"`
	Sb          null.Int    `json:"sb"`
	Cs          null.Int    `json:"cs"`
	Bb          null.Int    `json:"bb"`
	So          null.Int    `json:"so"`
	Sow         null.Float  `json:"sow"`
	Ba          null.Float  `json:"ba"`
	Obp         null.Float  `json:"obp"`
	Slg         null.Float  `json:"slg"`
	Ops         null.Float  `json:"ops"`
	Tb          null.Int    `json:"tb"`
	Gdp         null.Int    `json:"gdp"`
	Hbp         null.Int    `json:"hbp"`
	Sh          null.Int    `json:"sh"`
	Sf          null.Int    `json:"sf"`
	Ibb         null.Int    `json:"ibb"`
	Roe         null.Int    `json:"roe"`
	Babip       null.Float  `json:"babip"`
	Topsplus    null.Int    `json:"topsplus"`
	Sopsplus    null.Int    `json:"sopsplus"`
	Createddate null.Time   `json:"createddate"`
}

// Baserunner represents data for a baserunning
type Baserunner struct {
	ID          int         `json:"id"`
	Teamabbrev  string      `json:"teamabbrev"`
	Name        null.String `json:"name"`
	Age         null.Int    `json:"age"`
	Pa          null.Int    `json:"pa"`
	Roe         null.Int    `json:"roe"`
	Xi          null.Int    `json:"xi"`
	Rspct       null.String `json:"rspct"`
	Sbo         null.Int    `json:"sbo"`
	Sb          null.Int    `json:"sb"`
	Cs          null.Int    `json:"cs"`
	Sbpct       null.String `json:"sbpct"`
	Sb2         null.Int    `json:"sb2"`
	Cs2         null.Int    `json:"cs2"`
	Sb3         null.Int    `json:"sb3"`
	Cs3         null.Int    `json:"cs3"`
	Sbh         null.Int    `json:"sbh"`
	Csh         null.Int    `json:"csh"`
	Po          null.Int    `json:"po"`
	Pcs         null.Int    `json:"pcs"`
	Oob         null.Int    `json:"oob"`
	Oob1        null.Int    `json:"oob1"`
	Oob2        null.Int    `json:"oob2"`
	Oob3        null.Int    `json:"oob3"`
	Oobhm       null.Int    `json:"oobhm"`
	Bt          null.Int    `json:"bt"`
	Xbtpct      null.String `json:"xbtpct"`
	Firsts      null.Int    `json:"firsts"`
	Firsts2     null.Int    `json:"firsts2"`
	Firsts3     null.Int    `json:"firsts3"`
	Firstd      null.Int    `json:"firstd"`
	Firstd3     null.Int    `json:"firstd3"`
	Firstdh     null.Int    `json:"firstdh"`
	Seconds     null.Int    `json:"seconds"`
	Seconds3    null.Int    `json:"seconds3"`
	Secondsh    null.Int    `json:"secondsh"`
	Createddate null.Time   `json:"createddate"`
}
