package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// create a minimal JSON payload with a given next url
func makePayload(next string) string {
    return fmt.Sprintf(`{"count":1,"next":"%s","previous":null,"results":[{"name":"foo","url":"/foo"}]}`,
        next)
}

func TestCommandMapPagination(t *testing.T) {
    // helper to capture stdout from a function
    captureOutput := func(f func()) string {
        var buf bytes.Buffer
        orig := os.Stdout
        r, w, _ := os.Pipe()
        os.Stdout = w
        f()
        w.Close()
        io.Copy(&buf, r)
        os.Stdout = orig
        return buf.String()
    }
    // prepare an httptest server that serves two pages
    calls := 0
    var ts *httptest.Server
    ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        calls++
        switch calls {
        case 1:
            // first call returns a next link
            fmt.Fprintln(w, makePayload(ts.URL+"/page2"))
        case 2:
            // second call returns no next
            fmt.Fprintln(w, makePayload(""))
        default:
            // treat any additional calls as a fresh start (page1)
            fmt.Fprintln(w, makePayload(ts.URL+"/page2"))
        }
    }))
    defer ts.Close()

    // override baseURL and clear state
    baseURL = ts.URL + "/page1"
    resetMapState()
    if interactionManager.InteractionCount != 0 {
        t.Fatalf("resetMapState left interactions: %d", interactionManager.InteractionCount)
    }

    // first invocation should use the first page URL
    out := captureOutput(func() {
        if err := commandMap(); err != nil {
            t.Fatalf("first call failed: %v", err)
        }
    })
    if currentURL != ts.URL+"/page2" {
        t.Errorf("expected currentURL to be next page, got %q", currentURL)
    }
    if interactionManager.InteractionCount != 1 {
        t.Errorf("expected 1 interaction recorded, got %d", interactionManager.InteractionCount)
    }
    if interactionManager.Interactions[1].Next != ts.URL+"/page2" {
        t.Errorf("recorded next URL incorrect: %q", interactionManager.Interactions[1].Next)
    }
    // output should include the area name on its own line
    if !strings.Contains(strings.TrimSpace(out), "foo") {
        t.Errorf("output did not contain expected name:\n%s", out)
    }

    // second invocation should use the next link and then reset
    if err := commandMap(); err != nil {
        t.Fatalf("second call failed: %v", err)
    }
    if currentURL != "" {
        t.Errorf("expected currentURL to reset after reaching end, got %q", currentURL)
    }
    if interactionManager.InteractionCount != 2 {
        t.Errorf("expected 2 interactions recorded, got %d", interactionManager.InteractionCount)
    }
    if interactionManager.Interactions[2].Next != "" {
        t.Errorf("second recorded next URL incorrect: %q", interactionManager.Interactions[2].Next)
    }

    // a third call should start again at baseURL
    if err := commandMap(); err != nil {
        t.Fatalf("third call failed: %v", err)
    }
    if currentURL != ts.URL+"/page2" {
        t.Errorf("after resetting, expected to fetch baseURL and then next, got %q", currentURL)
    }
    if interactionManager.InteractionCount != 3 {
        t.Errorf("expected 3 interactions recorded, got %d", interactionManager.InteractionCount)
    }
}
