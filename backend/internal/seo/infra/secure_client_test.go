package infra

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecureClient_Hardcore(t *testing.T) {
	oldValidator := ipValidator
	defer func() { ipValidator = oldValidator }()

	secureClient := CreateSecureClient()
	fetcher := NewWebFetcher(secureClient)

	t.Run("Security: SSRF Protection", func(t *testing.T) {
		ipValidator = oldValidator

		forbiddenIPs := []string{
			"http://127.0.0.1:8080",
			"http://192.168.1.1",
			"http://10.0.0.5",
		}

		for _, target := range forbiddenIPs {
			t.Run(target, func(t *testing.T) {
				u, _ := url.Parse(target)
				_, err := fetcher.Scan(context.Background(), u)

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "access to restricted IP denied")
			})
		}
	})

	t.Run("Resilience: Request Timeout", func(t *testing.T) {
		ipValidator = func(ip net.IP) bool { return true }

		slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer slowServer.Close()

		u, _ := url.Parse(slowServer.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := fetcher.Scan(ctx, u)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("Parsing: Malformed HTML & Limits", func(t *testing.T) {
		ipValidator = func(ip net.IP) bool { return true }

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<html><title>Correct Title</title><body><h1>Header 1</h1>`)
			fmt.Fprint(w, string(make([]byte, 1024*600)))
			fmt.Fprint(w, `<h1>Hidden Header</h1></body></html>`)
		}))
		defer ts.Close()

		u, _ := url.Parse(ts.URL)
		report, err := fetcher.Scan(context.Background(), u)

		require.NoError(t, err)
		assert.Equal(t, "Correct Title", report.Metadata.Title)
		assert.Contains(t, report.Metadata.H1, "Header 1")
		for _, h := range report.Metadata.H1 {
			assert.NotEqual(t, "Hidden Header", h)
		}
	})

	t.Run("Network: User-Agent Spoofing", func(t *testing.T) {
		ipValidator = func(ip net.IP) bool { return true }

		uaChan := make(chan string, 1)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uaChan <- r.Header.Get("User-Agent")
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		u, _ := url.Parse(ts.URL)
		_, _ = fetcher.Scan(context.Background(), u)

		assert.Equal(t, "SiteInspector/1.0", <-uaChan)
	})
}
