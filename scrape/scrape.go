package scrape

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/hnucamendi/ws-colly_lambda/urls"
)

type Scrape struct {
	TopHeadlines []*TopHeadlines
	URLS         map[string]string

	wg *sync.WaitGroup
}

type TopHeadlines struct {
	NewsSites []*NewsSite
}

type NewsSite struct {
	Headline string
}

func NewScrape() *Scrape {
	s := &Scrape{}
	s.URLS = urls.URLS
	return s
}

func (s *Scrape) ScrapeTopHeadLines(c *colly.Collector) error {
	c.OnHTML(".stack_condensed", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildAttr(".container__title container_lead-package__title container__title--emphatic hover container__title--emphatic-size-l2", "h2"))
		// s.TopHeadlines = append(s.TopHeadlines, &TopHeadlines{
		// 	NewsSites: []*NewsSite{
		// 		{
		// 			Headline: e.Text,
		// 		},
		// 	},
		// })
	})

	c.Visit(s.URLS["cnn"])

	// for _, v := range s.URLS {
	// 	s.wg.Add(1)
	// 	go func(v string) {
	// 		defer s.wg.Done()
	// 		c.Visit(v)
	// 		fmt.Println(v)
	// 	}(v)
	// }

	return nil
}
