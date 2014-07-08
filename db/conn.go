package db

import (
	"net/url"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type DB struct {
	cfg  *config
	pool *redis.Pool
}

type config struct {
	user     string
	password string
	db       uint8
	addr     string
}

type Conn interface {
	Get() redis.Conn
}

func parseDSN(dsn string) (*config, error) {
	cfg := new(config)

	if dsn == "" {
		dsn = "redis://:@localhost:6379/0"
	}

	u, err := url.Parse(dsn)

	if err != nil {
		return nil, err
	}

	if pass, ok := u.User.Password(); ok {
		cfg.password = pass
	}

	db := u.Path

	if len(db) > 1 && db[0] == '/' {
		db = db[1:len(db)]
	}

	idb, err := strconv.ParseUint(db, 10, 8)

	if err != nil {
		idb = 0
	}

	cfg.db = uint8(idb)
	cfg.addr = u.Host
	return cfg, nil
}

func Open(dataSourceName string) (Conn, error) {
	var err error

	db := new(DB)
	db.cfg, err = parseDSN(dataSourceName)

	if err != nil {
		return nil, err
	}

	db.pool = &redis.Pool{
		MaxIdle:     128,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			return db.dial()
		},
		TestOnBorrow: nil,
	}
	return db, nil
}

func (db *DB) Get() redis.Conn {
	return db.pool.Get()
}

func (db *DB) dial() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", db.cfg.addr)

	if err != nil {
		return nil, err
	}

	//if LogLevel >= 2 {
	//	conn = redis.NewLoggingConn(conn, Logger, "")
	//}

	if db.cfg.password != "" {
		if _, err := conn.Do("AUTH", db.cfg.password); err != nil {
			//Logln("h: invalid redis password")
			conn.Close()
			return nil, err
		}
	}

	if db.cfg.db != 0 {
		if _, err := conn.Do("SELECT", db.cfg.db); err != nil {
			//Logln("h: invalid redis password")
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
