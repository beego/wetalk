// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utils

import (
	"database/sql"
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/setting"
)

var ErrSphinxDBClosed = fmt.Errorf("SphinxDB already closed and give back to pools")

var SphinxPools *sphinxPools

type SphinxDB struct {
	alive bool
	pools *sphinxPools
	db    *sql.DB
	orm   orm.Ormer
}

type SphinxMeta struct {
	Total      int
	TotalFound int64
	Time       float32
}

func (s *SphinxDB) RawQuery(query string, scans ...interface{}) (int64, error) {
	if !s.alive {
		return 0, ErrSphinxDBClosed
	}
	return s.orm.Raw(query).QueryRows(scans...)
}

func (s *SphinxDB) RawValuesFlat(query string, list *orm.ParamsList, column string) (int64, error) {
	if !s.alive {
		return 0, ErrSphinxDBClosed
	}
	return s.orm.Raw(query).ValuesFlat(list, column)
}

func (s *SphinxDB) ShowMeta() (*SphinxMeta, error) {
	if !s.alive {
		return nil, ErrSphinxDBClosed
	}
	meta := SphinxMeta{}
	if _, err := s.orm.Raw("SHOW META").RowsToStruct(&meta, "Variable_name", "Value"); err != nil {
		return nil, err
	}
	return &meta, nil
}

func (s *SphinxDB) Close() {
	if !s.alive {
		return
	}
	s.alive = false
	s.pools.giveBackDB(s)
}

func (s *SphinxDB) ping() error {
	return s.db.Ping()
}

func (s *SphinxDB) close() {
	s.db.Close()
}

type sphinxPools struct {
	conns int
	pools chan *SphinxDB
}

func (s *sphinxPools) GetConn() (sdb *SphinxDB, err error) {
	select {
	case sdb = <-s.pools:
	default:
	}

	if sdb != nil {
		if sdb.ping() == nil {
			sdb.alive = true
			return
		}

		sdb.close()
	}

	sdb, err = s.newConn()
	if sdb != nil {
		sdb.alive = true
	}
	return
}

func (s *sphinxPools) giveBackDB(sdb *SphinxDB) {
	sdb = &SphinxDB{
		alive: false,
		pools: s,
		db:    sdb.db,
		orm:   sdb.orm,
	}

	select {
	case s.pools <- sdb:
	default:
		sdb.close()
	}
}

func (s *sphinxPools) newConn() (*SphinxDB, error) {
	var (
		err error
		db  *sql.DB
		o   orm.Ormer
	)

	if db, err = sql.Open("sphinx", "root:root@tcp("+setting.SphinxHost+")/?loc=UTC"); err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(2)

	o, err = orm.NewOrmWithDB("sphinx", "sphinx", db)
	if err != nil {
		return nil, err
	}

	sdb := &SphinxDB{
		alive: true,
		pools: s,
		db:    db,
		orm:   o,
	}
	return sdb, nil
}

func (s *sphinxPools) close() {
	close(s.pools)

	for p := range s.pools {
		p.db.Close()
	}
}

func InitSphinxPools() error {
	if SphinxPools != nil {
		SphinxPools.close()
	} else {
		SphinxPools = &sphinxPools{
			conns: setting.SphinxMaxConn,
			pools: make(chan *SphinxDB, setting.SphinxMaxConn),
		}
	}

	for i := 0; i < setting.SphinxMaxConn; i++ {
		if sdb, err := SphinxPools.newConn(); err != nil {
			return err
		} else {
			SphinxPools.giveBackDB(sdb)
		}
	}

	return nil
}
