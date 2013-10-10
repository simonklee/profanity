package wordfilter

import (
	"testing"
)

var smallList = []string{"fuck", "duck", "puck", "suck", "eff"}
var largeList = []string{"@", "@A", "@AB", "A", "AB", "ABC", "B", "BC", "BCD", "C", "CD", "CDE", "D", "DE", "DEF", "E", "EF", "EFG", "F", "FG", "FGH", "G", "GH", "GHI", "H", "HI", "HIJ", "I", "IJ", "IJK", "J", "JK", "JKL", "K", "KL", "KLM", "L", "LM", "LMN", "M", "MN", "MNO", "N", "NO", "NOP", "O", "OP", "OPQ", "P", "PQ", "PQR", "Q", "QR", "QRS", "R", "RS", "RST", "S", "ST", "STU", "T", "TU", "TUV", "U", "UV", "UVW", "V", "VW", "VWX", "W", "WX", "WXY", "X", "XY", "XYZ", "Y", "YZ", "YZB", "Z", "ZA", "ZAB"}

type ProfanityTest struct {
	in, out string
}

func TestStringReplacer(t *testing.T) {
	tests := []*ProfanityTest{
		{"foo", "foo"},
		{"foo fuck", "foo ****"},
		{"foo fUCK", "foo ****"},
		{"foo uck", "foo uck"},
		{"foo ffuck", "foo f****"},
		{"eff", "***"},
	}

	repl := NewStringReplacer()
	repl.Reload(smallList)

	for i, x := range tests {
		if out := repl.Replace(x.in); out != x.out {
			t.Fatalf("#%d: expected %s, got %s", i, x.out, out)
		}
	}
}

func TestSetReplacer(t *testing.T) {
	tests := []*ProfanityTest{
		{"foo", "foo"},
		{"foo fuck", "foo ****"},
		{"foo fUCk", "foo ****"},
		{"foo uck", "foo uck"},
		{"foo ffuck", "foo ffuck"},
		{"eff", "***"},
		{"eff\n", "***\n"},
	}

	repl := NewSetReplacer()
	repl.Reload(smallList)

	for i, x := range tests {
		if out := repl.Replace(x.in); out != x.out {
			t.Fatalf("#%d: expected %s, got %s", i, x.out, out)
		}
	}
}

func BenchmarkBoyer(b *testing.B) {
	repl := NewStringReplacer()
	repl.Reload(largeList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repl.Replace("EFG")
	}
}

func BenchmarkSmallBoyerList(b *testing.B) {
	repl := NewStringReplacer()
	repl.Reload(smallList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repl.Replace("foo fuck")
	}
}

func BenchmarkSet(b *testing.B) {
	repl := NewSetReplacer()
	repl.Reload(largeList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repl.Replace("EFG")
	}
}

func BenchmarkSmallSetList(b *testing.B) {
	repl := NewSetReplacer()
	repl.Reload(smallList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repl.Replace("foo fuck")
	}
}

func BenchmarkLargeInputString(b *testing.B) {
	repl := NewStringReplacer()
	repl.Reload(smallList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repl.Replace("foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck")
	}
}

func BenchmarkLargeInputSet(b *testing.B) {
	repl := NewSetReplacer()
	repl.Reload(smallList)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repl.Replace("foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck foo fuck")
	}
}
