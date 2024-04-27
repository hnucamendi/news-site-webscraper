// lambda function to scrape news websites for headlines
package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gocolly/colly/v2"
	"github.com/hnucamendi/ws-colly_lambda/scfg"
)

type NewsSite struct {
	Site        string
	URL         string
	TopHeadline []*TopHeadline
}

type TopHeadline struct {
	Title       string
	Description string
	AritcleURL  string
	ImageURL    string
}

func sendToSQS(jb []byte, s *session.Session) error {
	ssmSvc := ssm.New(s)
	param, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("/lambda/prod/ws-colly/ws-colly-lambda-sqs-url"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	queueURL := *param.Parameter.Value

	sqsSvc := sqs.New(s)

	_, err = sqsSvc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(jb)),
		QueueUrl:     &queueURL,
	})

	return err
}

func scrapeTopHeadLines(c *colly.Collector, cfg *scfg.ScrapeConfig) ([]byte, error) {
	newsSites := &NewsSite{}
	c.OnHTML(cfg.Containers.TopHeadlinesContainer, func(e *colly.HTMLElement) {
		title := e.ChildText(cfg.TitleQuery)
		description := e.ChildText(cfg.DescriptionQuery)
		articleURL := e.ChildAttr(cfg.URLQuery, "href")
		imageURL := e.ChildAttr(cfg.ImageURLQuery, "src")

		if cfg.URLChopped {
			if s := e.ChildAttr(cfg.URLQuery, "href")[0]; s == '/' {
				articleURL = fmt.Sprintf("%s%s", cfg.URLPrefix, e.ChildAttr(cfg.URLQuery, "href"))
			}
		}

		newsSites.TopHeadline = append(newsSites.TopHeadline, &TopHeadline{
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
		newsSites.Site = strings.Split(r.URL.Host, ".")[1]
		newsSites.URL = r.URL.String()
	})

	c.Visit(cfg.URL)
	c.Wait()

	bytes, err := json.Marshal(newsSites)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func HandleRequest() error {
	c := colly.NewCollector(colly.Async(true), colly.UserAgent("ws-colly"))
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	json, err := scrapeTopHeadLines(c, scfg.CNNConfig())
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = sendToSQS(json, sess)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func main() {
	lambda.Start(HandleRequest)
}
