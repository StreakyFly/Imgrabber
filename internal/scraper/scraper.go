package scraper

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ScrapeImages(url string) ([]string, error) {
	htmlContent, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}

	return parseImages(htmlContent, url)
}

func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL %s: HTTP status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of URL %s: %w", url, err)
	}

	//fmt.Println(string(body))
	return string(body), nil
}

func parseImages(htmlContent string, baseURL string) ([]string, error) {
	var imageURLs []string
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return imageURLs, nil // end of content
		case html.StartTagToken,
			html.SelfClosingTagToken:
			token := tokenizer.Token()
			//fmt.Println(token.Data)

			if token.Data == "img" {
				for _, attr := range token.Attr {
					if attr.Key == "src" {
						absoluteURL, err := resolveURL(baseURL, attr.Val)
						if err != nil {
							return nil, err
						}
						imageURLs = append(imageURLs, absoluteURL)
					}
				}
			}
		default:
			//fmt.Println("tokenType:", tokenType, "data:", tokenizer.Token().Data)
			continue
		}
	}
}

func resolveURL(baseURL, relativeURL string) (string, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	parsedURL, err := parsedBaseURL.Parse(relativeURL)
	if err != nil {
		return "", fmt.Errorf("failed to resolve URL: %w", err)
	}

	return parsedURL.String(), nil
}
