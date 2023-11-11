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
		e.ChildText(".container__field-links")
		s.NewsSite.TopHeadlines = append(s.NewsSite.TopHeadlines, &TopHeadlines{
			NewsArticle: []*NewsArticle{
				{
					Title:       e.ChildText(".container__title_url-text"),
					Description: e.ChildText(".container__headline-text"),
					AritcleURL:  fmt.Sprintf("https://us.cnn.com%s", e.ChildAttr("a[href]", "href")),
				},
			},
		})
	})

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

	return nil
}
''