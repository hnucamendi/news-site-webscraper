package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gocolly/colly/v2"
	"github.com/hnucamendi/ws-colly_lambda/scrape"
)

func sendToSQS(json string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	ssmSvc := ssm.New(sess)
	param, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("/lambda/prod/ws-colly/ws-colly-lambda-sqs-url"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	queueURL := *param.Parameter.Value

	sqsSvc := sqs.New(sess)

	_, err = sqsSvc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(json),
		QueueUrl:     &queueURL,
	})

	return err
}

func HandleRequest() (string, error) {
	c := colly.NewCollector(colly.Async(true), colly.UserAgent("ws-colly"))
	s := scrape.NewScrape()

	cfg := &map[string]*scrape.ScrapeConfig{
		"cnn": {
			TitleQuery:       ".container__title_url-text",
			DescriptionQuery: ".container__headline-text",
			PaginationQuery:  "",
			URLQuery:         "a[href]",
			ImageURLQuery:    "img[src]",
			URL:              scrape.URLs["cnn"],
			URLPrefix:        "https://us.cnn.com",
			URLChopped:       true,
			Pagination:       false,
			WaitForLoad:      false,
			Containers: &scrape.SiteConfigContainer{
				TopHeadlinesContainer: ".zone__items",
			},
		},
	}

	json, err := s.ScrapeTopHeadLines(c, *cfg)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	err = sendToSQS(json)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return json, nil

}

func main() {
	lambda.Start(HandleRequest)
}
