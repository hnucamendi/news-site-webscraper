package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func uploadToDynamo() error {
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

	resp, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &queueURL,
		AttributeNames: []*string{
			aws.String("All"),
		},
	})
	if err != nil {
		return err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	for _, message := range resp.Messages {
		// Prepare item for DynamoDB
		item := map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(*message.MessageId),
			},
			"Body": {
				S: aws.String(*message.Body),
			},
			"TimeStamp": {
				S: aws.String(*message.Attributes["SentTimestamp"]),
			},
		}

		// Put item in DynamoDB
		input := &dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String("ws-colly-dynamo-table"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			return err
		}
	}

	return err
}

func HandleRequest() error {
	err := uploadToDynamo()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
