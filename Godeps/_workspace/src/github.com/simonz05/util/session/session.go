package session

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/simonz05/util/kvstore"
)

const (
	FullMask uint8 = (1 << iota)
	AdminMask
)

type Session struct {
	Id        string `json:"-"`
	Mask      uint8  `json:"m,omitempty"`
	ProfileID int    `json:"user_id,omitempty"`
}

type sessionString struct {
	ProfileID string `json:"user_id,omitempty"`
}

func (p *Session) HasAdmin() bool {
	return (p.Mask&AdminMask) != 0 && p.ProfileID != 0
}

func (p *Session) HasFull() bool {
	return (p.Mask & FullMask) != 0
}

func (p *Session) Set(mask uint8) {
	p.Mask |= mask
}

func (p *Session) Unset(mask uint8) {
	p.Mask &^= mask
}

var (
	// lifetime of session in seconds
	DefaultLifetime = 24 * 60 * 60
)

type Storage interface {
	Read(string) (*Session, error)
	Write(*Session) error
}

type redisBackend struct {
	prefix     string
	lifetime   int
	persistent bool
	db         *kvstore.KVStore
}

func NewRedisBackend(dns, prefix string, persistent bool) (Storage, error) {
	db, err := kvstore.Open(dns)

	if err != nil {
		return nil, err
	}

	conn := db.Get()
	_, err = conn.Do("Ping")
	conn.Close()

	if err != nil {
		return nil, err
	}

	return &redisBackend{
		prefix:     prefix,
		db:         db,
		lifetime:   DefaultLifetime,
		persistent: persistent,
	}, nil
}

func (rs *redisBackend) key(id string) string {
	return fmt.Sprintf("%s:%s", rs.prefix, id)
}

func (rs *redisBackend) Read(id string) (*Session, error) {
	conn := rs.db.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", rs.key(id)))

	if err != nil {
		return nil, err
	}

	s := new(Session)

	if err = json.Unmarshal(data, s); err != nil {
		// try to read profileID as string
		ss := new(sessionString)

		if err = json.Unmarshal(data, ss); err != nil {
			return nil, err
		}

		profileID, err := strconv.Atoi(ss.ProfileID)

		if err != nil {
			return nil, err
		}

		s.ProfileID = profileID
	}

	s.Id = id
	return s, nil
}

func (rs *redisBackend) Write(s *Session) error {
	data, err := json.Marshal(s)

	if err != nil {
		return err
	}

	conn := rs.db.Get()
	defer conn.Close()

	if rs.persistent {
		_, err = conn.Do("SET", rs.key(s.Id), data)
	} else {
		_, err = conn.Do("SETEX", rs.key(s.Id), rs.lifetime, data)
	}
	return err
}
