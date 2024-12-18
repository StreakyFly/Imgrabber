package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	DownloadDir = "./downloads"
	MaxWorkers  = 16
)

func DownloadImages(urls []string) error {
	err := os.MkdirAll(DownloadDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create a directory for images: %w", err)
	}

	urlCh := make(chan string)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < MaxWorkers; i++ {
		wg.Add(1)
		go worker(urlCh, &wg)
	}

	// Send URLs to the channel and close it
	go func() {
		for _, url := range urls {
			urlCh <- url
		}
		close(urlCh)
	}()

	wg.Wait() // wait for all workers to finish
	fmt.Println("All downloads complete.")
	return nil
}

func worker(urlCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urlCh {
		if err := downloadImage(url); err != nil {
			fmt.Printf("Failed to download %s: %v\n", url, err)
		}
	}
}

func downloadImage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status for URL %s: %d", url, resp.StatusCode)
	}

	uniqueName := uniqueFileName(url)
	filename := filepath.Join(DownloadDir, uniqueName)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving file %s: %w", filename, err)
	}

	fmt.Printf("Downloaded '%s': %s\n", uniqueName, url)
	return nil
}

func uniqueFileName(url string) string {
	urlParts := strings.Split(url, "/")
	nameParts := strings.Split(urlParts[len(urlParts)-1], ".")
	filename, extension := nameParts[0], nameParts[1]
	return fmt.Sprintf("%s_%d.%s", filename, time.Now().UnixNano(), extension)
}
