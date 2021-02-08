package db

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDB(cfg *config.Config) (*sqlx.DB, error) {
	datasource := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable", 
		cfg.DBUser, 
		cfg.DBPasswd, 
		cfg.DBName,
	)
	db, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		log.Fatalln(err)
	}
	return db, err
}
