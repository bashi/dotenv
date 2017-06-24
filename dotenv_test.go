package dotenv

import (
	"os"
	"testing"
)

type ParsedLine struct {
	key   string
	value string
}

func TestParseLine(t *testing.T) {
	if err := os.Setenv("TEST_ENV", "test_env"); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		"FOO=bar",
		"FOO=bar=baz",
		"FOO=./a/b/c",
		"FOO=/a/b/c",
		"BAR=$TEST_ENV",
		"BAR=$TEST_ENV/a/b/c",
	}
	expects := []ParsedLine{
		ParsedLine{"FOO", "bar"},
		ParsedLine{"FOO", "bar=baz"},
		ParsedLine{"FOO", "./a/b/c"},
		ParsedLine{"FOO", "/a/b/c"},
		ParsedLine{"BAR", "test_env"},
		ParsedLine{"BAR", "test_env/a/b/c"},
	}

	for i, line := range lines {
		expected := expects[i]
		key, value, err := parseLine(line, "")
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
