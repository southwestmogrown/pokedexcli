package locations

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExploreRequiresArgument(t *testing.T) {
    err := Explore([]string{})
    if err == nil {
        t.Fatal("expected usage error when explore has no arguments")
    }
    if !strings.Contains(err.Error(), "usage: explore") {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestExploreByNameAndCache(t *testing.T) {
    calls := 0
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        calls++
        if r.URL.Path != "/location-area/forest" {
            t.Fatalf("unexpected path: %s", r.URL.Path)
        }
        fmt.Fprintln(w, `{"pokemon_encounters":[{"pokemon":{"name":"pikachu","url":"/pikachu"}},{"pokemon":{"name":"oddish","url":"/oddish"}}]}`)
    }))
    defer ts.Close()

    baseURL = ts.URL + "/location-area"
    Reset()
    resetRequestCache()

    out1 := captureOutput(func() {
        if err := Explore([]string{"forest"}); err != nil {
            t.Fatalf("first explore failed: %v", err)
        }
    })

    if !strings.Contains(out1, "Exploring forest...") {
        t.Fatalf("missing explore intro, got: %s", out1)
    }
    if !strings.Contains(out1, "Found 2 Pokemon:") {
        t.Fatalf("missing pokemon header, got: %s", out1)
    }
    if !strings.Contains(out1, "pikachu") || !strings.Contains(out1, "oddish") {
        t.Fatalf("expected pokemon names in output, got: %s", out1)
    }
    if !strings.Contains(out1, "- pikachu") || !strings.Contains(out1, "- oddish") {
        t.Fatalf("expected pokemon bullet list output, got: %s", out1)
    }
    if strings.Contains(out1, "*** CACHE HIT ***") {
        t.Fatalf("first explore should not be cache hit: %s", out1)
    }
    if !strings.Contains(out1, "*** PAGE COMPLETE ***") {
        t.Fatalf("missing completion message: %s", out1)
    }

    out2 := captureOutput(func() {
        if err := Explore([]string{"forest"}); err != nil {
            t.Fatalf("second explore failed: %v", err)
        }
    })

    if !strings.Contains(out2, "*** CACHE HIT ***") {
        t.Fatalf("expected cache hit on second explore: %s", out2)
    }
    if calls != 1 {
        t.Fatalf("expected one HTTP request due to caching, got %d", calls)
    }
}

func TestExploreByID(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/location-area/12" {
            t.Fatalf("expected id path, got: %s", r.URL.Path)
        }
        fmt.Fprintln(w, `{"pokemon_encounters":[{"pokemon":{"name":"zubat","url":"/zubat"}}]}`)
    }))
    defer ts.Close()

    baseURL = ts.URL + "/location-area"
    Reset()
    resetRequestCache()

    out := captureOutput(func() {
        if err := Explore([]string{"12"}); err != nil {
            t.Fatalf("explore by id failed: %v", err)
        }
    })

    if !strings.Contains(out, "zubat") {
        t.Fatalf("expected zubat in output, got: %s", out)
    }
    if !strings.Contains(out, "Exploring 12...") {
        t.Fatalf("missing explore intro for id, got: %s", out)
    }
    if !strings.Contains(out, "Found 1 Pokemon:") {
        t.Fatalf("missing pokemon header for id, got: %s", out)
    }
    if !strings.Contains(out, "- zubat") {
        t.Fatalf("expected pokemon bullet list output for id, got: %s", out)
    }
}
