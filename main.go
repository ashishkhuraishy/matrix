package main

func main() {
	baseUrl := "https://go.dev"

	// Scrape(baseUrl)

	NewCrawler(baseUrl)

	select {}
}
