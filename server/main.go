// tServer project main.go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vyegres/tServer/server/storage"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func main() {
	db, e := sql.Open("mysql", "root@/scan?charset=utf8")
	defer db.Close()
	if e != nil {
		fmt.Println("Не удалось подключится к БД: ", e)
	}

	bs := new(storage.BlockStorage)
	bs.DB = db
	is := new(storage.IpStorage)
	is.DB = db

	rpc.Register(bs)
	rpc.Register(is)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}
}
