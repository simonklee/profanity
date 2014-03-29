package wordlist

import (
	"fmt"

	"github.com/simonz05/profanity/db"
	"github.com/simonz05/profanity/third_party/github.com/garyburd/redigo/redis"
	"github.com/simonz05/profanity/util"
)

// RedisWordlist is a redis backed wordlist implementation.
type RedisWordlist struct {
	lang string
	key  string
	conn db.Conn
}

func NewRedisWordlist(conn db.Conn, lang string) *RedisWordlist {
	return &RedisWordlist{
		lang: lang,
		key:  fmt.Sprintf("profanity:wordlist:%s", lang),
		conn: conn,
	}
}

func (w *RedisWordlist) Count() (int, error) {
	conn := w.conn.Get()
	defer conn.Close()
	return redis.Int(conn.Do("ZCARD", w.key))
}

func (w *RedisWordlist) Get(count, offset int) ([]string, error) {
	conn := w.conn.Get()
	defer conn.Close()
	starting_offset := util.IntMax(offset, 0)
	ending_offset := util.IntMax((starting_offset+count)-1, 1)
	return redis.Strings(conn.Do("ZRANGE", w.key, starting_offset, ending_offset))
}

func (w *RedisWordlist) Set(words []string) error {
	conn := w.conn.Get()
	defer conn.Close()
	conn.Send("MULTI")

	for i := 0; i < len(words); i++ {
		conn.Send("ZADD", w.key, 0, words[i])
	}
	_, err := conn.Do("EXEC")
	return err
}

func (w *RedisWordlist) Delete(words []string) error {
	conn := w.conn.Get()
	defer conn.Close()
	conn.Send("MULTI")

	for i := 0; i < len(words); i++ {
		conn.Send("ZREM", w.key, words[i])
	}

	_, err := conn.Do("EXEC")
	return err
}

func (w *RedisWordlist) Replace(words []string) error {
	if err := w.Empty(); err != nil {
		return err
	}
	return w.Set(words)
}

func (w *RedisWordlist) Empty() error {
	conn := w.conn.Get()
	_, err := conn.Do("DEL", w.key)
	conn.Close()
	return err
}
