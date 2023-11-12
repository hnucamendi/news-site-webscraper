package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/hnucamendi/ws-colly_lambda/scrape"
)

func main() {
	c := colly.NewCollector(colly.Async(true), colly.UserAgent("ws-colly"))
	s := scrape.NewScrape()

	if err := s.ScrapeTopHeadLines(c, scrape.CNNConfig()); err != nil {
		fmt.Println(err)
	}
}
