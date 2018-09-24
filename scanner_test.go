package toki

import (
	"testing"
)

const (
	NUMBER Token = iota + 1
	PLUS
	STRING
	COMMA
)

func TestScanner(t *testing.T) {
	input := "1  + 2+3 + happy birthday  "
	t.Logf("input: %s", input)
	s := NewScanner(
		[]Def{
			{Token: NUMBER, Pattern: "[0-9]+"},
			{Token: PLUS, Pattern: `\+`},
			{Token: STRING, Pattern: "[a-z]+"},
		})
	s.SetInput(input)
	expected := []Token{
		NUMBER,
		PLUS,
		NUMBER,
		PLUS,
		NUMBER,
		PLUS,
		STRING,
		STRING,
		EOF,
	}
	for _, e := range expected {
		r := s.Next()
		if e != r.Token {
			t.Fatalf("expected %v, got %v", e, r.Token)
		} else {
			t.Log(r)
		}
	}
}

func TestTrimWhitespace(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		trimWhitespace bool
		expectedValues []string
		def            []Def
	}{
		{
			name:           "trim whitespace (default)",
			input:          " a , b , c ",
			expectedValues: []string{"a", ",", "b", ",", "c"},
			trimWhitespace: true,
			def: []Def{
				{Token: STRING, Pattern: `[a-z]+`},
				{Token: COMMA, Pattern: ","},
			},
		},
		{
			// Note that, in this case, whitespace handling responsibility is
			// moved from the internal `skip` method, to instead be a function
			// of the STRING token pattern, increasing flexability of whitespace
			// handling.
			name:           "do not trim whitespace",
			input:          " a , b , c ",
			expectedValues: []string{" a ", ",", " b ", ",", " c "},
			trimWhitespace: false,
			def: []Def{
				{Token: STRING, Pattern: `([a-z]|\s)+`},
				{Token: COMMA, Pattern: ","},
			},
		},
	}

	for _, testCase := range testCases {
		test := func(t *testing.T) {
			s := NewScanner(testCase.def)
			s.SetInput(testCase.input)
			s.SetTrimWhitespace(testCase.trimWhitespace)
			if s.GetTrimWhitespace() != testCase.trimWhitespace {
				t.Log("GetTrimWhiteSpace returned the wrong value.")
				t.Error("Expected: ", testCase.trimWhitespace, "Got: ", s.GetTrimWhitespace())
			}
			for _, expectedValue := range testCase.expectedValues {
				r := s.Next()
				if r.Token == EOF {
					break
				}
				if expectedValue != string(r.Value) {
					t.Errorf("Expected: '%v' Got: '%v'", expectedValue, string(r.Value))
				}
			}
		}
		t.Run(testCase.name, test)
	}
}
