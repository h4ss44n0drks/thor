package logdb

import (
	"database/sql"
	"fmt"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/vechain/thor/thor"
	"sync"
)

//FilterOption option filter
type FilterOption struct {
	FromBlock uint32
	ToBlock   uint32
	Address   thor.Address // always a contract address
	Topics    [5]*thor.Hash
}

//LDB manages all logs
type LDB struct {
	path          string
	db            *sql.DB
	sqliteVersion string
	m             sync.RWMutex
}

//OpenDB open a logdb
func OpenDB(path string) (*LDB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	s, _, _ := sqlite3.Version()
	if err != nil {
		return nil, err
	}
	return &LDB{
		path:          path,
		db:            db,
		sqliteVersion: s,
	}, nil
}

//Insert insert logs into db
func (db *LDB) Insert(logs []*Log) error {
	if len(logs) == 0 {
		return nil
	}
	db.m.Lock()
	defer db.m.Unlock()

	stmt := ""
	for _, log := range logs {
		stmt += "insert into log(blockID ,blockNumber ,txID ,txOrigin ,address ,data ,topic0 ,topic1 ,topic2 ,topic3 ,topic4) values " + fmt.Sprintf(" ('%v',%v,'%v','%v','%v','%s','%v','%v','%v','%v','%v'); ", log.blockID.String(),
			log.blockNumber,
			log.txID,
			log.txOrigin,
			log.address,
			string(log.data),
			formatHash(log.topic0),
			formatHash(log.topic1),
			formatHash(log.topic2),
			formatHash(log.topic3),
			formatHash(log.topic4))
	}
	return db.ExecInTransaction(stmt)
}

//Filter return logs with options
func (db *LDB) Filter(options []*FilterOption) ([]*Log, error) {
	stmt := "select * from log where 1 "
	if len(options) != 0 {
		for _, op := range options {
			stmt += fmt.Sprintf(" or ( blockNumber >= %v and blockNumber <= %v and address = %v and topic0 = %v and topic1 = %v and topic2 = %v and topic3 = %v and topic4 = %v ) ",
				op.FromBlock,
				op.ToBlock,
				op.Address,
				formatHash(op.Topics[0]),
				formatHash(op.Topics[1]),
				formatHash(op.Topics[2]),
				formatHash(op.Topics[3]),
				formatHash(op.Topics[4]))
		}
	}
	return db.Query(stmt)
}

//ExecInTransaction execute sql in a transaction
func (db *LDB) ExecInTransaction(sqlStmt string, args ...interface{}) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(sqlStmt, args...); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//Query query logs
func (db *LDB) Query(stmt string) ([]*Log, error) {
	db.m.RLock()
	defer db.m.RUnlock()
	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*Log
	for rows.Next() {
		dbLog := &DBLog{}
		err = rows.Scan(
			&dbLog.blockID,
			&dbLog.blockNumber,
			&dbLog.txID,
			&dbLog.txOrigin,
			&dbLog.address,
			&dbLog.data,
			&dbLog.topic0,
			&dbLog.topic1,
			&dbLog.topic2,
			&dbLog.topic3,
			&dbLog.topic4)
		if err != nil {
			return nil, err
		}

		log, err := dbLog.toLog()
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return logs, nil
}

//Path return db's directory
func (db *LDB) Path() string {
	return db.path
}