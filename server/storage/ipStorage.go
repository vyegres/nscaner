package storage

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type IpStorage struct {
	DB *sql.DB
}

func (is *IpStorage) AddIP(ip uint32, reply *int) error {
	_, e := is.DB.Exec("INSERT INTO IP (ip) VALUES (?)", ip)
	if e != nil {
		return e
	}
	return nil
}
