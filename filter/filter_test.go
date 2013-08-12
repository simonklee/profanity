package filter

import (
	"testing"
)

var smallList = []string{
	"fuck",
	"duck",
	"puck",
	"suck",
	"eff",
}

var largeList = []string{
	"A", "AB", "ABC", "B", "BC", "BCD", "C", "CD", "CDE", "D", "DE", "DEF", "E", "EF", "EFG", "F", "FG", "FGH", "G", "GH", "GHI", "H", "HI", "HIJ", "I", "IJ", "IJK", "J", "JK", "JKL", "K", "KL", "KLM", "L", "LM", "LMN", "M", "MN", "MNO", "N", "NO", "NOP", "O", "OP", "OPQ", "P", "PQ", "PQR", "Q", "QR", "QRS", "R", "RS", "RST", "S", "ST", "STU", "T", "TU", "TUV", "U", "UV", "UVW", "V", "VW", "VWX", "W", "WX", "WXY", "X", "XY", "XYZ", "Y", "YZ", "YZB", "Z", "ZA", "ZAB", "B", "BC", "BCD", "C", "CD", "CDE", "D", "DE", "DEF", "E", "EF", "EFG", "F", "FG", "FGH", "G", "GH", "GHI", "H", "HI", "HIJ", "I", "IJ", "IJK", "J", "JK", "JKL", "K", "KL", "KLM", "L", "LM", "LMN", "M", "MN", "MNO", "N", "NO", "NOP", "O", "OP", "OPQ", "P", "PQ", "PQR", "Q", "QR", "QRS", "R", "RS", "RST", "S", "ST", "STU", "T", "TU", "TUV", "U", "UV", "UVW", "V", "VW", "VWX", "W", "WX", "WXY", "X", "XY", "XYZ", "Y", "YZ", "YZB", "Z", "ZA", "ZAB", "@", "@A", "@AB", "C", "CD", "CDE", "D", "DE", "DEF", "E", "EF", "EFG", "F", "FG", "FGH", "G", "GH", "GHI", "H", "HI", "HIJ", "I", "IJ", "IJK", "J", "JK", "JKL", "K", "KL", "KLM", "L", "LM", "LMN", "M", "MN", "MNO", "N", "NO", "NOP", "O", "OP", "OPQ", "P", "PQ", "PQR", "Q", "QR", "QRS", "R", "RS", "RST", "S", "ST", "STU", "T", "TU", "TUV", "U", "UV", "UVW", "V", "VW", "VWX", "W", "WX", "WXY", "X", "XY", "XYZ", "Y", "YZ", "YZB", "Z", "ZA", "ZAB", "@", "@A", "@AB", "@", "@A", "@AB", "D", "DE", "DEF", "E", "EF", "EFG", "F", "FG", "FGH", "G", "GH", "GHI", "H", "HI", "HIJ", "I", "IJ", "IJK", "J", "JK", "JKL", "K", "KL", "KLM", "L", "LM", "LMN", "M", "MN", "MNO", "N", "NO", "NOP", "O", "OP", "OPQ", "P", "PQ", "PQR", "Q", "QR", "QRS", "R", "RS", "RST", "S", "ST", "STU", "T", "TU", "TUV", "U", "UV", "UVW", "V", "VW", "VWX", "W", "WX", "WXY", "X", "XY", "XYZ", "Y", "YZ", "YZB", "Z", "ZA", "ZAB",
}

type ProfanityTest struct {
	in, out string
}

func initFilter() {
}

func TestFilter(t *testing.T) {
	pfilter := NewFilter()
	pfilter.Replace(smallList)

	tests := []*ProfanityTest{
		{"foo", "foo"},
		{"foo fuck", "foo ****"},
		{"foo uck", "foo uck"},
		{"foo ffuck", "foo f****"},
		{"eff", "***"},
	}

	for i, x := range tests {
		if out := pfilter.Sanitize(x.in); out != x.out {
			t.Fatalf("#%d: expected %s, got %s", i, x.out, out)
		}
	}
}

func TestBlacklist(t *testing.T) {
	pfilter := NewFilter()
	pfilter.Replace(smallList)

	if len(pfilter.Blacklist()) != len(smallList) {
		t.Fatalf("expected %d got %d", len(smallList), len(pfilter.Blacklist()))
	}

	newItems := []string{
		"mother",
	}

	pfilter.Replace(newItems)

	if len(pfilter.Blacklist()) != 1 {
		t.Fatalf("expected 1 got %d", len(pfilter.Blacklist()))
	}

	pfilter.Update(smallList)

	if len(pfilter.Blacklist()) != len(smallList)+1 {
		t.Fatalf("expected %d got %d, %v", len(smallList)+1, len(pfilter.Blacklist()), pfilter.Blacklist())
	}

	newItems = []string{
		"mother",
	}

	oldlen := len(pfilter.Blacklist())
	pfilter.Update(newItems)

	if len(pfilter.Blacklist()) != oldlen {
		t.Fatalf("expected %d got %d", oldlen, len(pfilter.Blacklist()))
	}
}

// func generateBlacklist(n int) []string {
// 	blacklist := make([]string, 0, n)
// 	for b := 'A'; b <= 'Z' ; b++ {
// 		for i:=0; i < 26; i++ {
// 			char := int(b)+i
//
// 			for j := 0; j < 3; j++ {
// 				w := ""
// 				for k := 0; k <= j; k++ {
// 					c := char+k
// 					if c > 'Z' {
// 						c = 'A'+k-1
// 					}
// 					w = fmt.Sprintf("%s%c", w, c)
// 				}
// 				fmt.Printf("\"%s\", ", w)
// 			}
// 		}
// 	}
// 	return blacklist
// }

func BenchmarkBoyer(b *testing.B) {
	pfilter := NewFilter()
	pfilter.Replace(largeList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pfilter.Sanitize("EFG")
	}
}

func BenchmarkSmallList(b *testing.B) {
	pfilter := NewFilter()
	pfilter.Replace(smallList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pfilter.Sanitize("foo fuck")
	}
}
