package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hnucamendi/ws-colly_lambda/urls"
	"golang.org/x/net/html"
)

type SiteConfigContainer struct {
	TopHeadlinesContainer string
}

type ScrapeConfig struct {
	TestQuery        func(...string) string
	TitleQuery       string
	DescriptionQuery string
	URLQuery         string
	ImageURLQuery    string
	URL              string
	URLPrefix        string
	URLChopped       bool
	Pagination       bool
	PaginationQuery  string
	Containers       *SiteConfigContainer
}

type Scrape struct {
	NewsSite NewsSite
	URLS     map[string]string
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
	ImageURL    string
}

func CNNConfig() *ScrapeConfig {
	sc := &ScrapeConfig{
		TitleQuery:       ".container__title_url-text",
		DescriptionQuery: ".container__headline-text",
		PaginationQuery:  "",
		URLQuery:         "a[href]",
		ImageURLQuery:    "img[src]",
		URL:              urls.URLS["cnn"],
		URLPrefix:        "https://us.cnn.com",
		URLChopped:       true,
		Pagination:       false,
		Containers: &SiteConfigContainer{
			TopHeadlinesContainer: ".zone__items",
		},
	}

	return sc
}

func ViceConfig() *ScrapeConfig {
	sc := &ScrapeConfig{
		TitleQuery:       ".vice-card__content",
		DescriptionQuery: ".vice-card-dek",
		PaginationQuery:  ".loading-lockup-infinite__button",
		URLQuery:         "a[href]",
		ImageURLQuery:    "picture[source]",
		URL:              urls.URLS["vice"],
		URLChopped:       false,
		URLPrefix:        "",
		Pagination:       false,
		Containers: &SiteConfigContainer{
			// TopHeadlinesContainer: ".container",
			TopHeadlinesContainer: "body",
		},
	}

	return sc
}

func NewScrape() *Scrape {
	s := &Scrape{}
	s.URLS = urls.URLS
	return s
}

func (s *Scrape) ScrapeTopHeadLines(c *colly.Collector, cfg *ScrapeConfig) error {
	c.OnHTML(cfg.Containers.TopHeadlinesContainer, func(e *colly.HTMLElement) {
		// fmt.Println(e)
		w := io.Writer(os.Stdout)
		if err := html.Render(w, e.DOM.Nodes[0]); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(w)

		title := e.ChildText(cfg.TitleQuery)
		description := e.ChildText(cfg.DescriptionQuery)
		articleURL := e.ChildAttr(cfg.URLQuery, "href")
		imageURL := e.ChildAttr(cfg.ImageURLQuery, "src")

		if cfg.URLChopped {
			if s := e.ChildAttr(cfg.URLQuery, "href")[0]; s == '/' {
				articleURL = fmt.Sprintf("%s%s", cfg.URLPrefix, e.ChildAttr(cfg.URLQuery, "href"))
			}
		}

		s.NewsSite.TopHeadlines = append(s.NewsSite.TopHeadlines, &TopHeadlines{
			NewsArticle: []*NewsArticle{
				{
					Title:       title,
					Description: description,
					AritcleURL:  articleURL,
					ImageURL:    imageURL,
				},
			},
		})
	})

	if cfg.Pagination {
		c.OnHTML(cfg.PaginationQuery, func(h *colly.HTMLElement) {
			t := h.ChildAttr("a", "href")
			c.Visit(t)
		})
	}

	c.OnRequest(func(r *colly.Request) {
		s.NewsSite.Site = strings.Split(r.URL.Host, ".")[1]
		s.NewsSite.URL = r.URL.String()
	})

	c.Visit(cfg.URL)
	c.Wait()

	_, err := json.Marshal(s.NewsSite)
	if err != nil {
		return err
	}

	// fmt.Println(string(bytes))

	return nil
}
