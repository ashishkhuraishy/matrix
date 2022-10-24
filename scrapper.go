package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func Scrape(baseUrl string) {
	urls := requestData(baseUrl)
	fmt.Println(urls)
}

func requestData(baseUrl string) []string {
	// log.Println(baseUrl)
	var urls []string

	baseUrl = strings.TrimSuffix(baseUrl, "/")

	resp, err := http.Get(baseUrl)
	if err != nil {
		log.Println(err)
		return urls
	}

	token := html.NewTokenizer(resp.Body)
	for {
		tt := token.Next()
		if tt == html.ErrorToken {
			if token.Err() == io.EOF {
				break
			}

			log.Printf("Error: %v\n", token.Err())
			break
		}

		tag, hasAttr := token.TagName()
		// fmt.Println(string(token.Text()))
		if string(tag) == "a" && hasAttr {
			var key, value []byte
			hasMore := true
			for hasMore {
				key, value, hasMore = token.TagAttr()
				if string(key) == "href" {
					uri := string(value)
					if !strings.HasPrefix(uri, "/") {
						continue
					}

					val := fmt.Sprintf("%v%v", baseUrl, uri)
					u, err := url.Parse(val)
					if err != nil {
						log.Println(err)
						continue
					}

					val = fmt.Sprintf("%v://%v%v", u.Scheme, u.Host, uri)

					// fmt.Println(val)
					urls = append(urls, val)
				}
			}
		}

		// if hasAttr {
		// 	for {
		// 		key, value, hasMore := token.TagAttr()
		// 		fmt.Printf("%v: %v - %v\n", string(tag), string(key), string(value))
		// 		if !hasMore {
		// 			break
		// 		}
		// 	}
		// }
	}

	return urls

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return

	// }

	// fmt.Println(string(body))
}
