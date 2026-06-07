package infra_test

import (
	"backend/internal/seo/infra"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecureClient_Hardcore(t *testing.T) {
	oldValidator := infra.IpValidator
	defer func() { infra.IpValidator = oldValidator }()

	secureClient := infra.CreateSecureClient()
	fetcher := infra.NewWebFetcher(secureClient)

	t.Run("Security: SSRF Protection", func(t *testing.T) {
		infra.IpValidator = oldValidator

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
		infra.IpValidator = func(ip net.IP) bool { return true }

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

	t.Run("Network: User-Agent Spoofing", func(t *testing.T) {
		infra.IpValidator = func(ip net.IP) bool { return true }

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
