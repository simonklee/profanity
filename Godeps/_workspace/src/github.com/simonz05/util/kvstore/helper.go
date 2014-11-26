package kvstore

import (
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

// Ints is a helper that converts a multi-bulk command reply to a []int.
// If err is not equal to nil, then Ints returns nil, err.  If one if the
// multi-bulk items is not a bulk value or nil, then Ints returns an error.
func Ints(reply interface{}, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		result := make([]int, 0, len(reply))
		n := 0
		for _, v := range reply {
			if v == nil {
				continue
			}
			p, ok := v.([]byte)
			if !ok {
				return nil, fmt.Errorf("kvstore: unexpected element type for Ints, got type %T", v)
			}

			nr, err := strconv.ParseInt(string(p), 10, 0)

			if err != nil {
				return nil, err
			}

			n++
			result = append(result, int(nr))
		}
		return result[:n], nil
	case nil:
		return nil, redis.ErrNil
	case redis.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("kvstore: unexpected type for Ints, got type %T", reply)
}

// Bytes is a helper that converts a multi-bulk command reply to a [][]byte.
// If err is not equal to nil, then Ints returns nil, err.  If one if the
// multi-bulk items is not a bulk value or nil, then Ints returns an error.
func Bytes(reply interface{}, err error) ([][]byte, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		result := make([][]byte, 0, len(reply))
		n := 0
		for _, v := range reply {
			var b []byte
			switch raw := v.(type) {
			case []byte:
				b = raw
			case string:
				b = []byte(raw)
			case nil:
				continue
			default:
				return nil, fmt.Errorf("kvstore: unexpected type for Bytes, got type %T", v)
			}

			n++
			result = append(result, b)
		}
		return result[:n], nil
	case nil:
		return nil, redis.ErrNil
	case redis.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("kvstore: unexpected type for Bytes, got type %T", reply)
}
