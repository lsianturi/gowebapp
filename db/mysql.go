package db

import (
	"database/sql"
	// how to comment this
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lsianturi/gowebapp/model"
)

// Config is bla
type Config struct {
	ConnectString string
}

// InitDb is bla
func InitDb(cfg Config) (*MySQLDb, error) {
	dbConn, err := sqlx.Connect("mysql", cfg.ConnectString)
	if err != nil {
		return nil, err
	}

	p := &MySQLDb{dbConn: dbConn}
	if err := p.dbConn.Ping(); err != nil {
		return nil, err
	}
	if err := p.createTableIfNotExist(); err != nil {
		return nil, err
	}
	if err := p.prepareSQLStatements(); err != nil {
		return nil, err
	}
	return p, nil
}

// MySQLDb is bla
type MySQLDb struct {
	dbConn *sqlx.DB

	sqlSelectPeople *sqlx.Stmt
	sqlInsertPerson *sqlx.NamedStmt
	sqlSelectPerson *sql.Stmt
}

func (p *MySQLDb) createTableIfNotExist() error {
	createSQL := `
		CREATE TABLE IF NOT EXISTS people (
			id int(11) NOT NULL AUTO_INCREMENT,
			first_name varchar(30) NOT NULL,
			last_name varchar(30) NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1
	`

	rows, err := p.dbConn.Query(createSQL)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (p *MySQLDb) prepareSQLStatements() (err error) {

	if p.sqlSelectPeople, err = p.dbConn.Preparex(
		"SELECT id, firstname, lastname FROM people",
	); err != nil {
		return err
	}

	if p.sqlInsertPerson, err = p.dbConn.PrepareNamed(
		"INSERT INTO people (firstname, lastname) " +
			"VALUES (:firstname, :lastname)",
	); err != nil {
		return err
	}

	if p.sqlSelectPerson, err = p.dbConn.Prepare(
		"SELECT id, firstname, lastname FROM people WHERE id = ?",
	); err != nil {
		return err
	}
	return nil
}

// SelectPeople is bla
func (p *MySQLDb) SelectPeople() ([]*model.Person, error) {
	people := make([]*model.Person, 0)
	if err := p.sqlSelectPeople.Select(&people); err != nil {
		return nil, err
	}
	return people, nil
}
