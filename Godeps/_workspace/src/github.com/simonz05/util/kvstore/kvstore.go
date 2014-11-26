package kvstore

import (
	"net/url"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/simonz05/util/log"
)

type KVStore struct {
	cfg  *config
	Pool *redis.Pool
}

type config struct {
	password string
	db       uint8
	addr     string
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
		db = db[1:]
	}

	idb, err := strconv.ParseUint(db, 10, 8)

	if err != nil {
		idb = 0
	}

	cfg.db = uint8(idb)
	cfg.addr = u.Host
	return cfg, nil
}

func Open(dataSourceName string) (*KVStore, error) {
	var err error

	kvstore := new(KVStore)
	kvstore.cfg, err = parseDSN(dataSourceName)

	if err != nil {
		return nil, err
	}

	kvstore.Pool = &redis.Pool{
		MaxIdle:     128,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			return kvstore.dial()
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return kvstore, nil
}

func (kvstore *KVStore) Get() redis.Conn {
	return kvstore.Pool.Get()
}

func (kvstore *KVStore) Close() error {
	return kvstore.Pool.Close()
}

func (kvstore *KVStore) dial() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", kvstore.cfg.addr)

	if err != nil {
		return nil, err
	}

	if kvstore.cfg.password != "" {
		if _, err := conn.Do("AUTH", kvstore.cfg.password); err != nil {
			log.Errorf("Redis AUTH err: %v", err)
			conn.Close()
			return nil, err
		}
	}

	if kvstore.cfg.db != 0 {
		if _, err := conn.Do("SELECT", kvstore.cfg.db); err != nil {
			log.Errorf("Redis SELECT err: %v", err)
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
