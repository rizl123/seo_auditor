package infra_test

import (
	"backend/internal/seo/infra"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebFetcher_Scan(t *testing.T) {
	tests := []struct {
		name            string
		html            string
		expectTitle     string
		expectDesc      string
		expectH1        int
		expectOGImage   string
		expectCanonical string
	}{
		{
			name: "valid full page",
			html: `
			<html>
				<head>
					<title>Go Test Page</title>
					<meta name="description" content="SEO testing is fun">
					<meta property="og:image" content="https://img.com/a.png">
					<link rel="canonical" href="https://example.com">
				</head>
				<body>
					<h1>Hello World</h1>
					<h1>Second Title</h1>
				</body>
			</html>`,
			expectTitle:     "Go Test Page",
			expectDesc:      "SEO testing is fun",
			expectH1:        2,
			expectOGImage:   "https://img.com/a.png",
			expectCanonical: "https://example.com",
		},
		{
			name: "missing meta tags",
			html: `
			<html>
				<head><title>Only Title</title></head>
				<body><h1>One</h1></body>
			</html>`,
			expectTitle: "Only Title",
			expectH1:    1,
		},
		{
			name: "broken whitespace",
			html: `
			<html>
				<head>
					<title>   spaced    title   </title>
				</head>
				<body>
					<h1>   hello    world   </h1>
				</body>
			</html>`,
			expectTitle: "spaced title",
			expectH1:    1,
		},
		{
			name: "empty page",
			html: `
			<html><head></head><body></body></html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.html))
			}))
			defer server.Close()

			fetcher := infra.NewWebFetcher(http.DefaultClient)

			targetURL, _ := url.Parse(server.URL)
			raw, err := fetcher.Scan(context.Background(), targetURL)

			assert.NoError(t, err)
			assert.NotNil(t, raw)

			if tt.expectTitle != "" {
				assert.Equal(t, tt.expectTitle, raw.Metadata.Title)
			}
			if tt.expectDesc != "" {
				assert.Equal(t, tt.expectDesc, raw.Metadata.Description)
			}
			if tt.expectH1 > 0 {
				assert.Len(t, raw.Metadata.H1, tt.expectH1)
			}
			if tt.expectOGImage != "" {
				assert.Equal(t, tt.expectOGImage, raw.Metadata.OgImage)
			}
			if tt.expectCanonical != "" {
				assert.Equal(t, tt.expectCanonical, raw.Metadata.Canonical)
			}
		})
	}
}

func TestWebFetcher_Scan_EmptyPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><head></head><body></body></html>`))
	}))
	defer server.Close()

	fetcher := infra.NewWebFetcher(http.DefaultClient)

	targetURL, _ := url.Parse(server.URL)
	raw, err := fetcher.Scan(context.Background(), targetURL)

	assert.NoError(t, err)
	assert.NotNil(t, raw)

	meta := raw.Metadata

	assert.Empty(t, meta.Title)
	assert.Empty(t, meta.Description)
	assert.Empty(t, meta.H1)
	assert.Empty(t, meta.Canonical)
	assert.Empty(t, meta.OgImage)
}

func TestWebFetcher_Scan_BrokenHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
			<html>
				<head>
					<title>Broken Title</title>
					<meta name="description" content="unclosed">
				<body>
					<h1>Still Works</h1>
		`))
	}))
	defer server.Close()

	fetcher := infra.NewWebFetcher(http.DefaultClient)

	targetURL, _ := url.Parse(server.URL)
	raw, err := fetcher.Scan(context.Background(), targetURL)

	require.NoError(t, err)
	require.NotNil(t, raw)
	require.NotNil(t, raw.Metadata)

	meta := raw.Metadata
	assert.Equal(t, "Broken Title", meta.Title)
	require.Len(t, meta.H1, 1)
	assert.Equal(t, "Still Works", meta.H1[0])
}

func TestWebFetcher_Scan_MultipleCanonical(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
			<html>
				<head>
					<link rel="canonical" href="https://first.com">
					<link rel="canonical" href="https://second.com">
				</head>
			</html>
		`))
	}))
	defer server.Close()

	fetcher := infra.NewWebFetcher(http.DefaultClient)

	targetURL, _ := url.Parse(server.URL)
	raw, err := fetcher.Scan(context.Background(), targetURL)

	assert.NoError(t, err)

	meta := raw.Metadata

	assert.Equal(t, "https://second.com", meta.Canonical)
}

func TestWebFetcher_Scan_WhitespaceNormalization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
			<html>
				<head>
					<title>
						Go     Test

						Page
					</title>
				</head>
				<body>
					<h1>
						Hello


						World
					</h1>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	fetcher := infra.NewWebFetcher(http.DefaultClient)

	targetURL, _ := url.Parse(server.URL)
	raw, err := fetcher.Scan(context.Background(), targetURL)

	assert.NoError(t, err)

	meta := raw.Metadata

	assert.Equal(t, "Go Test Page", meta.Title)
	assert.Equal(t, "Hello World", meta.H1[0])
}
