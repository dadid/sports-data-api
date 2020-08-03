package db

import (
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // provides SQL Server driver
	_ "github.com/go-sql-driver/mysql"   // provides MySQL driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // provides postgres driver
	"github.com/pkg/errors"
)

// Container holds the connection pool created with Config
type Container struct {
	Db   *sqlx.DB
	Conf Dbconfig
}

// Dbconfig holds database creds used for initializing a new Container
type Dbconfig struct {
	Rdbms    string `json:"rdbms"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	WinAuth  bool   `json:"winauth"`
}

var (
	errNilDb  = errors.New("sql.DB pointer is nil")
	errDbConf = errors.New("all database configuration fields must be set")
)

// Close performs the release of any resources that `sql/database` DB pool created.
func (dbc *Container) Close() (err error) {
	if dbc.Db == nil {
		err = errNilDb
		return
	}

	if err = dbc.Db.Close(); err != nil {
		err = errors.Wrap(err, "Error closing database connection.")
		return
	}
	return
}

// Open creates a database connection and pings the database to verify connection is alive
func (dbc *Container) Open() (err error) {
	// Check to see if any config string fields were left empty
	if dbc.Conf.Rdbms == "" || dbc.Conf.Host == "" || dbc.Conf.Port == "" || dbc.Conf.User == "" || dbc.Conf.Password == "" || dbc.Conf.Database == "" {
		err = errDbConf
		return
	}

	db, err := sqlx.Open(dbc.Conf.Rdbms, dbc.Conf.createConnString())
	if err != nil {
		err = errors.Wrap(err, "error opening connection to database")
		return
	}
	if err = db.Ping(); err != nil {
		err = errors.Wrapf(err, "error pinging %s", dbc.Conf.Host)
		return
	}
	dbc.Db = db
	return
}

func (conf *Dbconfig) createConnString() (connString string) {
	switch conf.Rdbms {
	case "postgres":
		connString = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", conf.Host, conf.Port, conf.Database, conf.User, conf.Password)
	case "mysql":
		connString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	case "sqlserver":
		if conf.WinAuth {
			connString = fmt.Sprintf("server=%s; port=%s; database=%s; trusted_connection=yes;", conf.Host, conf.Port, conf.Database)
		}
		connString = fmt.Sprintf("server=%s; port=%s; database=%s; user id=%s; password=%s", conf.Host, conf.Port, conf.Database, conf.User, conf.Password)
	}
	return
}
