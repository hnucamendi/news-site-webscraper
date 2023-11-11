package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	// s := scrape.NewScrape()

	// if err := s.ScrapeTopHeadLines(c); err != nil {
	// 	fmt.Println(err)
	// }

	//.container_lead-package .cnn .lazy
	c.OnHTML(".container", func(e *colly.HTMLElement) {
		t := map[string]string{
			"Title":       e.ChildText(".container__title_url-text"),
			"Description": e.ChildText(".container__headline-text"),
		}
		fmt.Println(t)
	})

	if err := c.Visit("https://www.cnn.com"); err != nil {
		fmt.Println(err)
	}

	// fmt.Println(s.TopHeadlines)
}
