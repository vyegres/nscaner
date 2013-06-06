package storage

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/go-sql-driver/mysql"
	"net"
)

type IpStorage struct {
	DB *sql.DB
}

func (is *IpStorage) AddIP(ipItem IpItem, reply *int) error {
	_, e := is.DB.Exec("INSERT INTO IP (ip, hello) VALUES (?,?)", ipItem.IP, ipItem.Hello)
	if e != nil {
		return e
	}
	return nil
}

type IpItem struct {
	IP    uint32
	Hello string
}

func (ipItem *IpItem) String() string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipItem.IP)
	ip := net.IP(ipByte)
	return ip.String()
}
