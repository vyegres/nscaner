package storage

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type BlockStorage struct {
	DB *sql.DB
}

func (bs *BlockStorage) GetNew(_, reply *uint32) error {
	var id uint32
	//Проверим нет ли блока с просроченной датой
	res := bs.DB.QueryRow("SELECT id FROM  `Block` WHERE timeDone=0 and (CURRENT_TIMESTAMP - timeCreate) > (10*60 * 60 *24)")
	e := res.Scan(&id)
	if e == nil {
		bs.DB.Exec("UPDATE Block set timeCreate=CURRENT_TIMESTAMP where id=?", id)
	} else {
		res := bs.DB.QueryRow("SELECT MAX(id) FROM Block")
		res.Scan(&id)
		id++
		bs.DB.Exec("INSERT into Block (id) VALUES  (?)", id)
	}
	*reply = id
	return nil
}

func (bs *BlockStorage) SetDone(id int, reply *int) error {
	_, e := bs.DB.Exec("UPDATE Block set timeDone=CURRENT_TIMESTAMP where id=?", id)
	if e != nil {
		return e
	}
	return nil
}
