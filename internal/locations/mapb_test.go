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

func captureOutputB(f func()) string {
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

func makePayloadB(prev string) string {
    return fmt.Sprintf(
        `{"count":1,"next":"","previous":"%s","results":[{"name":"bar","url":"/bar"}]}`,
        prev)
}

func TestMapBUsesPrevious(t *testing.T) {
    calls := 0
    var ts *httptest.Server
    ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        calls++
        switch calls {
        case 1:
            fmt.Fprintln(w, makePayloadB(ts.URL+"/page1"))
        case 2:
            fmt.Fprintln(w, makePayloadB(""))
        default:
            fmt.Fprintln(w, makePayloadB(""))
        }
    }))
    defer ts.Close()

    baseURLB = ts.URL + "/page2"
    ResetB()
    resetRequestCache()

    out := captureOutputB(func() {
        if err := MapB(); err != nil {
            t.Fatalf("first call failed: %v", err)
        }
    })
    if !strings.Contains(out, "*** REACHED BEGINNING ***") {
        t.Errorf("missing beginning log: %s", out)
    }
    if !strings.Contains(strings.TrimSpace(out), "bar") {
        t.Errorf("output missing name: %s", out)
    }
    if !strings.Contains(out, "Exploring location areas...") {
        t.Errorf("missing location intro: %s", out)
    }
    if !strings.Contains(out, "Found 1 Location Area:") {
        t.Errorf("missing location header: %s", out)
    }
    if !strings.Contains(out, "- bar") {
        t.Errorf("missing location bullet entry: %s", out)
    }
    if !strings.Contains(out, "*** PAGE COMPLETE ***") {
        t.Errorf("missing completion flair: %s", out)
    }
    if currentURL != ts.URL+"/page2" {
        t.Errorf("expected currentURL to be base page after first call, got %q", currentURL)
    }

    resetRequestCache()

    out2 := captureOutputB(func() {
        if err := MapB(); err != nil {
            t.Fatalf("second call failed: %v", err)
        }
    })
    if strings.Contains(out2, "*** REACHED BEGINNING ***") {
        t.Errorf("unexpected beginning log on second call: %s", out2)
    }
    if currentURL != ts.URL+"/page1" {
        t.Errorf("expected currentURL=%q, got %q", ts.URL+"/page1", currentURL)
    }

    resetRequestCache()
    out3 := captureOutputB(func() {
        if err := MapB(); err != nil {
            t.Fatalf("third call failed: %v", err)
        }
    })
    if !strings.Contains(out3, "*** REACHED BEGINNING ***") {
        t.Errorf("expected beginning log on third call: %s", out3)
    }
    if !strings.Contains(out3, "*** PAGE COMPLETE ***") {
        t.Errorf("missing completion flair on third call: %s", out3)
    }
    if currentURL != ts.URL+"/page2" {
        t.Errorf("expected currentURL back at base, got %q", currentURL)
    }
}
