package talking_clock

import "testing"

var testCases = []struct {
	input  string
	output string
}{
	{"23:59", "It's eleven fifty nine pm"},
	{"11:59", "It's eleven fifty nine am"},
	{"10:11", "It's ten eleven am"},
	{"15:19", "It's three nineteen pm"},
	{"01:00", "It's one am"},
	{"17:00", "It's five pm"},
	{"09:20", "It's nine twenty am"},
	{"13:30", "It's one thirty pm"},
	{"12:00", "It's twelve pm"},
	{"00:00", "It's twelve am"},
}

func TestTranslate(t *testing.T) {
	for _, test := range testCases {
		actual := Translate(test.input)
		if actual != test.output {
			t.Errorf("Expected \"%s\", got \"%s\"", test.output, actual)
		}
	}
}
