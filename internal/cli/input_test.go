package cli

import "testing"

func TestCleanInput(t *testing.T) {
    cases := []struct {
        input    string
        expected []string
    }{
        {input: " hello world ", expected: []string{"hello", "world"}},
        {input: "BulBaSauR  Reptar  MothrA", expected: []string{"bulbasaur", "reptar", "mothra"}},
    }

    for _, c := range cases {
        actual := cleanInput(c.input)
        if len(actual) != len(c.expected) {
            t.Fatalf("actual length %d does not match expected %d", len(actual), len(c.expected))
        }

        for i := range actual {
            if actual[i] != c.expected[i] {
                t.Fatalf("word %q does not match expected %q", actual[i], c.expected[i])
            }
        }
    }
}
