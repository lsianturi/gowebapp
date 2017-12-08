package db

import (
	"database/sql"

	"github.com/lsianturi/gowebapp/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	ConnectString string
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("mysql", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.createTableIfNotExist(); err != nil {
			return nil, err
		}
		if err := p.prepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

type pgDb struct {
	dbConn *sqlx.DB

	sqlSelectPeople *sqlx.Stmt
	sqlInsertPerson *sqlx.NamedStmt
	sqlSelectPerson *sqlx.Stmt
}

func (p *pgDb) createTableIfNotExist() error {
	create_sql := `
		CREATE TABLE IF NOT EXIST people (id AUTOINCREMENT NOT NULL PRIMARY KEY,
		first_name TEXT NOT NULL,
	    last_name TEXT NOT NULL);
	`

	if rows, err := p.dbConn.Query(create_sql); err != nil {
		return err
	} else {
		rows.Close()
	}
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {
	if p.sqlSelectPeople, err = p.dbConn.Preparex (
		"SELECT id, first_name, last_name FROM people",); err != nil {
			return err
		}
	}
	if p.sqlInsertPerson, err = p.dbConn.PrepareNamed (
		"INSERT INTO people (first_name, last_name) 
		   VALUES (:first, :last); SELECT LAST_INSERT_ID();",); err != nil {
		return err
	}
	if p.sqlSelectPerson, err = p.dbConn.Prepare(
		"SELECT id, first_name, last_name FROM people WHERE id = &1",
	); err != nil {
		return err
	}
	return nil
}

func (p *pgDb) ([]*model.Person, error) {
	people := make([]*model.Person, 0)
	if err := p.sqlSelectPeople.Select(&people); err != nil {
		return nil, err
	}
	return people, nil
}
