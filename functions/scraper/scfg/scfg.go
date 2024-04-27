package scfg

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
