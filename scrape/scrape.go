package scrape

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
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
	WaitForLoad      bool
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
	Title       string
	Description string
	AritcleURL  string
	ImageURL    string
}

var urls map[string]string = map[string]string{
	"google":             "https://www.google.com",
	"cnn":                "https://www.cnn.com",
	"nytimes":            "https://www.nytimes.com/",
	"nbcnews":            "https://www.nbcnews.com/",
	"cbsnews":            "https://www.cbsnews.com/",
	"vice":               "https://www.vice.com/en/section/news",
	"cnbc":               "https://www.cnbc.com/",
	"forbes":             "https://www.forbes.com/",
	"espn":               "https://www.espn.com/",
	"foxnews":            "https://www.foxnews.com/",
	"usatoday":           "https://www.usatoday.com/",
	"bbc":                "https://www.bbc.com/news/world/us_and_canada",
	"washingtonpost":     "https://www.washingtonpost.com/",
	"latimes":            "https://www.latimes.com/",
	"npr":                "https://www.npr.org/",
	"wsj":                "https://www.wsj.com/us-news",
	"theguardian":        "https://www.theguardian.com/us",
	"bloomberg":          "https://www.bloomberg.com/",
	"reuters":            "https://www.reuters.com/world/us/",
	"usnews":             "https://www.usnews.com/",
	"nationalgeographic": "https://www.nationalgeographic.com/",
	"apnews":             "https://www.apnews.com/us-news",
	"yahoo":              "https://news.yahoo.com/",
}

func Sleep(t int, d time.Duration) {
	time.Sleep(time.Duration(t) * d)
}

func CNNConfig() *ScrapeConfig {
	sc := &ScrapeConfig{
		TitleQuery:       ".container__title_url-text",
		DescriptionQuery: ".container__headline-text",
		PaginationQuery:  "",
		URLQuery:         "a[href]",
		ImageURLQuery:    "img[src]",
		URL:              urls["cnn"],
		URLPrefix:        "https://us.cnn.com",
		URLChopped:       true,
		Pagination:       false,
		WaitForLoad:      false,
		Containers: &SiteConfigContainer{
			TopHeadlinesContainer: ".zone__items",
		},
	}

	return sc
}

func ViceConfig() *ScrapeConfig {
	sc := &ScrapeConfig{
		TitleQuery:       ".latest-feed",
		DescriptionQuery: ".vice-card-dek",
		PaginationQuery:  ".loading-lockup-infinite__button",
		URLQuery:         "a[href]",
		ImageURLQuery:    "picture[source]",
		URL:              urls["vice"],
		URLChopped:       false,
		URLPrefix:        "",
		Pagination:       false,
		WaitForLoad:      true,
		Containers: &SiteConfigContainer{
			// TopHeadlinesContainer: ".container",
			TopHeadlinesContainer: "body",
		},
	}

	return sc
}

func NewScrape() *Scrape {
	s := &Scrape{}
	s.URLS = urls
	return s
}

func (s *Scrape) ScrapeTopHeadLines(c *colly.Collector, cfg *ScrapeConfig) error {
	c.OnHTML(cfg.Containers.TopHeadlinesContainer, func(e *colly.HTMLElement) {
		// w := io.Writer(os.Stdout)
		// if err := html.Render(w, e.DOM.Nodes[0]); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// fmt.Println(w)

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
			Title:       title,
			Description: description,
			AritcleURL:  articleURL,
			ImageURL:    imageURL,
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

	bytes, err := json.Marshal(s.NewsSite)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}
