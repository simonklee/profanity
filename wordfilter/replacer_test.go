package wordfilter

import (
	"testing"
)

var smallList = []string{"fuck", "duck", "puck", "suck", "eff"}
var largeList = []string{"@", "@A", "@AB", "A", "AB", "ABC", "B", "BC", "BCD", "C", "CD", "CDE", "D", "DE", "DEF", "E", "EF", "EFG", "F", "FG", "FGH", "G", "GH", "GHI", "H", "HI", "HIJ", "I", "IJ", "IJK", "J", "JK", "JKL", "K", "KL", "KLM", "L", "LM", "LMN", "M", "MN", "MNO", "N", "NO", "NOP", "O", "OP", "OPQ", "P", "PQ", "PQR", "Q", "QR", "QRS", "R", "RS", "RST", "S", "ST", "STU", "T", "TU", "TUV", "U", "UV", "UVW", "V", "VW", "VWX", "W", "WX", "WXY", "X", "XY", "XYZ", "Y", "YZ", "YZB", "Z", "ZA", "ZAB"}

type ProfanityTest struct {
	in, out string
}

func TestReplacer(t *testing.T) {
	pfilter := NewReplacer()
	pfilter.Reload(smallList)

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

func BenchmarkBoyer(b *testing.B) {
	pfilter := NewReplacer()
	pfilter.Reload(largeList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pfilter.Sanitize("EFG")
	}
}

func BenchmarkSmallList(b *testing.B) {
	pfilter := NewReplacer()
	pfilter.Reload(smallList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pfilter.Sanitize("foo fuck")
	}
}
