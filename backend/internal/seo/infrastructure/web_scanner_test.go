package infrastructure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebScanner_Scan(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
            <html>
                <head>
                    <title>Go Test Page</title>
                    <meta name="description" content="SEO testing is fun">
                </head>
                <body>
                    <h1>Hello World</h1>
                    <h1>Second Title</h1>
                </body>
            </html>
        `))
	}))
	defer server.Close()

	scanner := NewWebScanner(http.DefaultClient)

	url, _ := url.Parse(server.URL)
	report, err := scanner.Scan(context.Background(), url)

	assert.NoError(t, err)
	if assert.NotNil(t, report) {
		assert.Equal(t, "Go Test Page", report.Metadata.Title)
		assert.Equal(t, "SEO testing is fun", report.Metadata.Description)
		assert.Len(t, report.Metadata.H1, 2)
		assert.Equal(t, "Hello World", report.Metadata.H1[0])
	}
}

func TestWebScanner_Security(t *testing.T) {
	secureClient := CreateSecureClient()
	scanner := NewWebScanner(secureClient)

	t.Run("Should block local addresses", func(t *testing.T) {
		url, _ := url.Parse("http://127.0.0.1:8080")
		_, err := scanner.Scan(context.Background(), url)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "internal network access denied")
	})

	t.Run("Should block invalid protocol", func(t *testing.T) {
		url, _ := url.Parse("ftp://unsafe-site.com")
		_, err := scanner.Scan(context.Background(), url)

		assert.Error(t, err)
		assert.Equal(t, "invalid protocol", err.Error())
	})
}
