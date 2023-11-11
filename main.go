package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/hnucamendi/ws-colly_lambda/scrape"
)

func main() {
	c := colly.NewCollector(colly.Async(true))
	s := scrape.NewScrape()

	if err := s.ScrapeTopHeadLines(c); err != nil {
		fmt.Println(err)
	}

	// for i, v := range s.TopHeadlines {
	// 	fmt.Println(v.NewsSites[i])
	// }
}
