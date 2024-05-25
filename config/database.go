package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DbConnection interface {
	Conn() *sql.DB
}

type dbConnection struct {
	db *sql.DB
	cfg *Config
}

func (d *dbConnection) initDb() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.cfg.DbConfig.Host,
		d.cfg.DbConfig.Port,
		d.cfg.DbConfig.User,
		d.cfg.DbConfig.Password,
		d.cfg.DbConfig.Name,
	)
	db, err := sql.Open(d.cfg.DbConfig.Driver,dsn)
	if err != nil {
		return fmt.Errorf("Error Create Connection : %v",err)
	}
	checkConn := db.Ping()
	if checkConn != nil {
		return fmt.Errorf("Error Ping !! : %v",checkConn.Error())
	}
	d.db = db
	return nil
}

func (d *dbConnection) Conn() *sql.DB {
	return d.db
}

func NewDbConnection(config *Config) (DbConnection,error){
	conn := &dbConnection{
		cfg: config,
	}

	err := conn.initDb()
	if err != nil {
		return nil, err
	}
	return conn,nil
} 
