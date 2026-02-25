package repl

import (
	"fmt"
	"testing"
)

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
        fmt.Println(actual)
        if len(actual) != len(c.expected) {
            t.Errorf("Actual length: %d does not match Expected length: %d", len(actual), len(c.expected))
            t.Fail()
        }

        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]

            if word != expectedWord {
                t.Errorf("Word %s does not match expected word %s", word, expectedWord)
                t.Fail()
            }
        }
    }
}
