package infra

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebFetcher_Scan(t *testing.T) {
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

	fetcher := NewWebFetcher(http.DefaultClient)

	targetURL, _ := url.Parse(server.URL)
	report, err := fetcher.Scan(context.Background(), targetURL)

	assert.NoError(t, err)
	if assert.NotNil(t, report) {
		assert.Equal(t, "Go Test Page", report.Metadata.Title)
		assert.Equal(t, "SEO testing is fun", report.Metadata.Description)
		assert.Len(t, report.Metadata.H1, 2)
		assert.Equal(t, "Hello World", report.Metadata.H1[0])
	}
}
