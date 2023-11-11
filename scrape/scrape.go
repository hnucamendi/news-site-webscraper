package scrape

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/hnucamendi/ws-colly_lambda/urls"
)

type Scrape struct {
	NewsSite NewsSite
	URLS     map[string]string

	wg *sync.WaitGroup
}

type NewsSite struct {
	Site         string
	URL          string
	TopHeadlines []*TopHeadlines
}

type TopHeadlines struct {
	NewsArticle []*NewsArticle
}

type NewsArticle struct {
	Title       string
	Description string
	AritcleURL  string
}

func NewScrape() *Scrape {
	s := &Scrape{}
	s.URLS = urls.URLS
	return s
}

func (s *Scrape) ScrapeTopHeadLines(c *colly.Collector) error {
	c.OnHTML(".zone__items", func(e *colly.HTMLElement) {
		// e.ForEach(".container.container_lead-package.cnn", func(_ int, c *colly.HTMLElement) {
		e.ForEach(".container__title_url-text", func(_ int, t *colly.HTMLElement) {
			fmt.Printf("Title: %s\n", t.Text)
			e.ForEach(".container__headline-text", func(_ int, d *colly.HTMLElement) {
				fmt.Printf("Description: %s\n", d.Text)
				e.ForEach(".container__field-links", func(_ int, du *colly.HTMLElement) {
					// fmt.Printf("DESC: %s\n", d.Text)
					// fmt.Printf("Title: %s\nDescription: %s\nURL: %s\n", t.Text, d.Text, du.ChildAttr("a[href]", "href"))

					// s.NewsSite.TopHeadlines = append(s.NewsSite.TopHeadlines, &TopHeadlines{
					// 	NewsArticle: []*NewsArticle{
					// 		{
					// 			Title:       t.Text,
					// 			Description: d.Text,
					// 			AritcleURL:  fmt.Sprintf("https://us.cnn.com%s", du.ChildAttr("a[href]", "href")),
					// 		},
					// 	},
					// })
				})
			})
		})
	})
	// })

	c.OnRequest(func(r *colly.Request) {
		s.NewsSite.Site = strings.Split(r.URL.Host, ".")[1]
		s.NewsSite.URL = r.URL.String()
	})

	c.Visit(s.URLS["cnn"])
	c.Wait()

	bytes, err := json.Marshal(s.NewsSite)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

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
