package wordlist

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/simonz05/profanity/db"
	"github.com/simonz05/profanity/util"
)

var (
	largeList = []string{"@", "@A", "@AB", "A", "AB", "ABC", "B", "BC", "BCD", "C", "CD", "CDE", "D", "DE", "DEF", "E", "EF", "EFG", "F", "FG", "FGH", "G", "GH", "GHI", "H", "HI", "HIJ", "I", "IJ", "IJK", "J", "JK", "JKL", "K", "KL", "KLM", "L", "LM", "LMN", "M", "MN", "MNO", "N", "NO", "NOP", "O", "OP", "OPQ", "P", "PQ", "PQR", "Q", "QR", "QRS", "R", "RS", "RST", "S", "ST", "STU", "T", "TU", "TUV", "U", "UV", "UVW", "V", "VW", "VWX", "W", "WX", "WXY", "X", "XY", "XYZ", "Y", "YZ", "YZB", "Z", "ZA", "ZAB"}
	smallList = []string{"A", "AB", "ABC"}
	backends  []Wordlist
	once      sync.Once
)

type TestCase struct {
	name  string
	words []string
	list  Wordlist
}

func (t *TestCase) String() string {
	return fmt.Sprintf("%s - %s", reflect.TypeOf(t.list), t.name)
}

func initBackend() {
	c, _ := db.Open("redis://:@localhost:6379/15")
	backends = append(backends, NewRedisWordlist(c, "en_US"))
}

func testBackend(t *testing.T, test *TestCase) {
	list := test.list
	words := test.words

	if err := list.Empty(); err != nil {
		t.Fatalf("%s: expected err %v", test, err)
	}

	if cnt, err := list.Count(); cnt != 0 || err != nil {
		t.Fatalf("%s expected 0 got %d, err %v", test, cnt, err)
	}

	if values, err := list.Get(10, 0); len(values) != 0 || err != nil {
		t.Fatalf("%s expected 0 got %d, err %v", test, len(values), err)
	}

	if err := list.Set(words); err != nil {
		t.Fatalf("%s expected nil got %v", test, err)
	}

	if cnt, err := list.Count(); cnt != len(words) || err != nil {
		t.Fatalf("%s expected %d got %d, err %v", test, len(words), cnt, err)
	}

	expCnt := util.IntMin(10, len(words))

	if values, err := list.Get(10, 0); len(values) != expCnt || err != nil {
		t.Fatalf("%s expected %d got %d, err %v", test, expCnt, len(values), err)
	}

	if len(words) > 2 {
		if values, err := list.Get(1, 1); len(values) != 1 || err != nil {
			t.Fatalf("%s expected 1 got %d, err %v", test, len(values), err)
		}
	}

	if len(words) > 1 {
		if err := list.Delete(words[:1]); err != nil {
			t.Fatalf("%s expected nil got %v", test, err)
		}

		if cnt, err := list.Count(); cnt != len(words)-1 || err != nil {
			t.Fatalf("%s expected %d got %d, err %v", test, len(words)-1, cnt, err)
		}
	}

	if err := list.Empty(); err != nil {
		t.Fatalf("%s: expected err %v", test, err)
	}

	if cnt, err := list.Count(); cnt != 0 || err != nil {
		t.Fatalf("%s expected 0 got %d, err %v", test, cnt, err)
	}
}

func TestWordlist(t *testing.T) {
	once.Do(initBackend)

	tests := []TestCase{
		{name: "small list", words: smallList},
		{name: "large list", words: largeList},
	}

	for _, test := range tests {
		for _, backend := range backends {
			test.list = backend
			testBackend(t, &test)
		}
	}
}
