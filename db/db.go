package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb" // provides SQL Server driver
	_ "github.com/go-sql-driver/mysql"   // provides MySQL driver
	_ "github.com/lib/pq"                // provides postgres driver
	"github.com/pkg/errors"
)

// Container holds the connection pool created with Config
type Container struct {
	Db   *sql.DB
	conf config
}

// config holds connection string info used for instantiating a new DbContainer
type config struct {
	Rdbms    string `json:"rdbms"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	WinAuth  bool   `json:"winauth"`
}

var (
	// ErrNilDB - universal error for db is nil
	ErrNilDB = errors.New("sql.DB pointer is nil")
	// Dbconfig is the universal database configuration
	Dbconfig config
	_        = json.Unmarshal([]byte(os.Getenv("BR_API_DB_CONFIG")), &Dbconfig)
)

// Close performs the release of any resources that `sql/database` DB pool created.
func (dbc *Container) Close() (err error) {
	if dbc.Db == nil {
		err = ErrNilDB
		return
	}

	if err = dbc.Db.Close(); err != nil {
		err = errors.Wrap(err, "Error closing database connection.")
		return
	}
	return
}

// NewContainer initalizes a new Container object based on a config struct
func NewContainer(conf config) (dbc Container, err error) {
	// Check to see if any config string fields were left empty
	if conf.Rdbms == "" || conf.Host == "" || conf.Port == "" || conf.User == "" || conf.Password == "" || conf.Database == "" {
		err = errors.Errorf("All configuration fields must be set. %v", conf)
		return
	}
	dbc.conf = conf
	var sqlConnString string
	switch conf.Rdbms {
	case "postgres":
		sqlConnString = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", conf.Host, conf.Port, conf.Database, conf.User, conf.Password)
	case "mysql":
		sqlConnString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	case "sqlserver":
		if conf.WinAuth {
			sqlConnString = fmt.Sprintf("server=%s; port=%s; database=%s; trusted_connection=yes;", conf.Host, conf.Port, conf.Database)
		}
		sqlConnString = fmt.Sprintf("server=%s; port=%s; database=%s; user id=%s; password=%s", conf.Host, conf.Port, conf.Database, conf.User, conf.Password)
	}
	db, err := sql.Open(conf.Rdbms, sqlConnString)
	if err != nil {
		err = errors.Wrap(err, "error opening connection to database")
		return
	}
	// Ping verifies if the database connection is alive.
	if err = db.Ping(); err != nil {
		err = errors.Wrapf(err, "error pinging %s", conf.Host)
		return
	}

	dbc.Db = db
	return
}
