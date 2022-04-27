package strReplacer

import "testing"

func TestRandStr(t *testing.T) {
	type args struct {
		letter []string
		n      int
	}
	n := 20
	tests := []struct {
		name string
		args args
	}{
		{name: "number", args: args{letter: []string{"number"}, n: n}},
		{name: "lower", args: args{letter: []string{"lower"}, n: n}},
		{name: "upper", args: args{letter: []string{"upper"}, n: n}},
		{name: "special", args: args{letter: []string{"special"}, n: n}},
		{name: "lower_upper", args: args{letter: []string{"lower", "upper"}, n: n}},
		{name: "number_lower_upper", args: args{letter: []string{"lower", "upper", "number"}, n: n}},
		{name: "number_lower_upper_special", args: args{letter: []string{"lower", "upper", "number", "special"}, n: n}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := RandStr(tt.args.letter, tt.args.n)
			t.Log(s)

		})
	}
}
