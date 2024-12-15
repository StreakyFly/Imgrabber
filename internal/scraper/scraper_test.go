package scraper

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestScrapeImages(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    []string
		wantErr bool
	}{
		{"Invalid URL", "invalid-page", nil, true},
		{"URL with images", "valid-page", []string{"https://example.com/img1.jpg", "https://example.com/img2.jpg"}, false},
		{"URL with no images", "no-images", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock server with custom responses based on the URL
			handler := http.NewServeMux()
			handler.HandleFunc("/invalid-page", func(w http.ResponseWriter, r *http.Request) {
				http.NotFound(w, r)
			})
			handler.HandleFunc("/valid-page", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`
					<html>
						<body>
							<img src="https://example.com/img1.jpg"/>
							<img src="https://example.com/img2.jpg"/>
						</body>
					</html>
				`))
			})
			handler.HandleFunc("/no-images", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("<html><body>No images here.</body></html>"))
			})

			// Create a new test server
			ts := httptest.NewServer(handler)
			defer ts.Close()

			url := ts.URL + "/" + tt.url

			// Call the function under test
			got, err := ScrapeImages(url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScrapeImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScrapeImages() got = %v, want %v", got, tt.want)
			}
		})
	}
}
