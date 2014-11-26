package session

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/simonz05/util/assert"
)

type backendTest struct {
	got   *Session
	exp   *Session
	err   error
	sleep int
}

func TestBackend(t *testing.T) {
	ast := assert.NewAssert(t)

	DefaultLifetime = 1
	redisStorage, err := NewRedisBackend("redis://:@localhost:6379/15", "dev", false)
	ast.Nil(err)

	backends := []Storage{
		redisStorage,
	}

	backendTests := []backendTest{
		{
			got:   &Session{Id: "1"},
			exp:   &Session{Id: "1"},
			err:   nil,
			sleep: 0,
		},
		{
			got:   &Session{Id: "1", Mask: AdminMask | FullMask},
			exp:   &Session{Id: "1", Mask: AdminMask | FullMask},
			err:   nil,
			sleep: 0,
		},
		{
			got:   &Session{Id: "1"},
			exp:   nil,
			err:   nil,
			sleep: 2,
		},
	}

	for _, backend := range backends {
		for _, test := range backendTests {
			err := backend.Write(test.got)
			ast.Nil(err)

			if test.sleep > 0 {
				time.Sleep(time.Duration(int(time.Second) * test.sleep))
			}

			ses, err := backend.Read(test.got.Id)

			if test.err != nil {
				ast.Equal(test.err, err)
			} else {
				ast.Equal(test.exp, ses)
			}
		}
	}
}

func TestSession(t *testing.T) {
	ast := assert.NewAssert(t)

	p := &Session{}

	ast.True(!p.HasAdmin())
	ast.True(!p.HasFull())

	p.Set(AdminMask | FullMask)

	ast.True(!p.HasAdmin())
	ast.True(p.HasFull())

	p.ProfileID = 1

	ast.True(p.HasAdmin())
	ast.True(p.HasFull())

	p.Unset(AdminMask)

	ast.True(!p.HasAdmin())
	ast.True(p.HasFull())
}

func TestPersistance(t *testing.T) {
	ast := assert.NewAssert(t)

	p1 := &Session{}
	p1.Set(AdminMask)
	p1.ProfileID = 1

	buf, err := json.Marshal(p1)
	ast.Nil(err)

	p2 := &Session{}
	err = json.Unmarshal(buf, p2)

	//fmt.Printf("%s\n", buf)
	ast.Nil(err)
	ast.True(p2.HasAdmin())
	ast.True(!p2.HasFull())

	buf, err = json.Marshal(&Session{})

	ast.Nil(err)
	ast.Equal("{}", string(buf))
}

func BenchmarkPersistance(b *testing.B) {
	p1 := &Session{}
	p1.Set(AdminMask)
	p1.ProfileID = 1

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(p1)

		if err != nil {
			b.Fatal(err)
		}
	}
}
