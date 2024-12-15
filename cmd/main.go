package main

import (
	"Imgrabber/internal/scraper"
	"fmt"
	"log"
)

func main() {
	url := "YOUR_EPIC_WEBSITE_URL"
	imageURLs, err := scraper.ScrapeImages(url)
	if err != nil {
		log.Fatalf("Error scraping images: %v", err)
	}

	fmt.Println("Found image URLs:")
	for _, url := range imageURLs {
		fmt.Println(url)
	}
	fmt.Println("----------------------")
}
