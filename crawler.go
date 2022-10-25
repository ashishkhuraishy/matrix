package main

import (
	"log"
	"sync"
)

type Crawler struct {
	// base seed url
	BaseURL string

	// keeping track of all the
	// visiste urls
	visited map[string]bool

	urls   chan string
	parser chan []string

	// a mutex lock for adding
	mutex sync.RWMutex
}

func NewCrawler(baseUrl string) *Crawler {
	crawler := Crawler{
		BaseURL: baseUrl,
		visited: make(map[string]bool),
		urls:    make(chan string, 1),
		parser:  make(chan []string, 10),
	}

	go crawler.Crawl()
	go crawler.parse()
	crawler.urls <- baseUrl

	return &crawler
}

func (c *Crawler) Crawl() {
	for url := range c.urls {
		if c.check(url) {
			continue
		}

		go func(url string) {
			urls := requestData(url)
			log.Printf("found %d links from %v\n", len(urls), url)
			c.visit(url)

			c.parser <- urls
		}(url)
	}
}

func (c *Crawler) parse() {
	for urls := range c.parser {
		for _, url := range urls {
			if c.check(url) {
				continue
			}

			c.urls <- url
		}
	}
}

func (c *Crawler) visit(url string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.visited[url] = true
}

func (c *Crawler) check(url string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.visited[url]
}
