package cache

import (
	"testing"
	"time"
)

func TestReapLoopRemovesOldEntries(t *testing.T) {
    interval := 50 * time.Millisecond
    c := NewPokeCache(interval)

    c.Add("a", []byte{1})
    c.Add("b", []byte{2})

    if _, ok := c.Get("a"); !ok {
        t.Fatal("expected key a")
    }
    if _, ok := c.Get("b"); !ok {
        t.Fatal("expected key b")
    }

    time.Sleep(interval + 20*time.Millisecond)

    c.Add("c", []byte{3})

    if _, ok := c.Get("a"); ok {
        t.Error("expected key a to be reaped")
    }
    if _, ok := c.Get("b"); ok {
        t.Error("expected key b to be reaped")
    }
    if _, ok := c.Get("c"); !ok {
        t.Error("expected key c still present")
    }
}
