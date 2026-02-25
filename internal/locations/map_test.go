package locations

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

func captureOutput(f func()) string {
    var buf bytes.Buffer
    orig := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    f()
    w.Close()
    _, _ = io.Copy(&buf, r)
    os.Stdout = orig
    return buf.String()
}

func makePayload(next string) string {
    return fmt.Sprintf(`{"count":1,"next":"%s","previous":null,"results":[{"name":"foo","url":"/foo"}]}`,
        next)
}

func TestMapPagination(t *testing.T) {
    calls := 0
    var ts *httptest.Server
    ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        calls++
        switch calls {
        case 1:
            fmt.Fprintln(w, makePayload(ts.URL+"/page2"))
        case 2:
            fmt.Fprintln(w, makePayload(""))
        default:
            fmt.Fprintln(w, makePayload(ts.URL+"/page2"))
        }
    }))
    defer ts.Close()

    baseURL = ts.URL + "/page1"
    Reset()
    resetRequestCache()
    if interactionManager.InteractionCount != 0 {
        t.Fatalf("reset left interactions: %d", interactionManager.InteractionCount)
    }

    out := captureOutput(func() {
        if err := Map(); err != nil {
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
    if !strings.Contains(strings.TrimSpace(out), "foo") {
        t.Errorf("output missing name: %s", out)
    }
    if !strings.Contains(out, "Exploring location areas...") {
        t.Errorf("missing location intro: %s", out)
    }
    if !strings.Contains(out, "Found 1 Location Area:") {
        t.Errorf("missing location header: %s", out)
    }
    if !strings.Contains(out, "- foo") {
        t.Errorf("missing location bullet entry: %s", out)
    }
    if !strings.Contains(out, "*** PAGE COMPLETE ***") {
        t.Errorf("missing completion flair: %s", out)
    }

    resetRequestCache()
    if err := Map(); err != nil {
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

    resetRequestCache()
    out2 := captureOutput(func() {
        if err := Map(); err != nil {
            t.Fatalf("third call failed: %v", err)
        }
    })
    if !strings.Contains(out2, "*** WRAPPED TO START ***") {
        t.Errorf("wrap-around message missing: %s", out2)
    }
    if !strings.Contains(out2, "*** PAGE COMPLETE ***") {
        t.Errorf("missing completion flair on wrap call: %s", out2)
    }
    if currentURL != ts.URL+"/page2" {
        t.Errorf("after resetting, expected to fetch baseURL and then next, got %q", currentURL)
    }
    if interactionManager.InteractionCount != 3 {
        t.Errorf("expected 3 interactions recorded, got %d", interactionManager.InteractionCount)
    }

    currentURL = ts.URL + "/page1"
    lastFetchedURL = ""
    out3 := captureOutput(func() {
        if err := Map(); err != nil {
            t.Fatalf("cache call failed: %v", err)
        }
    })
    if !strings.Contains(out3, "*** CACHE HIT ***") {
        t.Errorf("expected cache hit on repeated URL, got: %s", out3)
    }
    if !strings.Contains(out3, "*** PAGE COMPLETE ***") {
        t.Errorf("missing completion flair on cache call: %s", out3)
    }
}
