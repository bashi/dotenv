package dotenv

import "testing"

type ParsedLine struct {
	key   string
	value string
}

func TestParseLine(t *testing.T) {
	lines := []string{
		"FOO=bar",
		"FOO=bar=baz",
	}
	expects := []ParsedLine{
		ParsedLine{"FOO", "bar"},
		ParsedLine{"FOO", "bar=baz"},
	}

	for i, line := range lines {
		expected := expects[i]
		key, value, err := parseLine(line)
		if err != nil {
			t.Fatalf("Failed to parse: %v\n", err)
			if key != expected.key {
				t.Fatalf("Key mismatch: %v != %v\n", expected.key, key)
			}
			if value != expected.value {
				t.Fatalf("Value mismatch: %v != %v\n", expected.value, value)
			}
		}
	}
}
