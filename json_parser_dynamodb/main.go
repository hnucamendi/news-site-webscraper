package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type client struct {
	sqs      *sqs.SQS
	dynamodb *dynamodb.DynamoDB
	ssm      *ssm.SSM
}

func initClient() *client {
	cl := &client{}

	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	cl.sqs = sqs.New(newSession)
	cl.dynamodb = dynamodb.New(newSession)
	cl.ssm = ssm.New(newSession)

	return cl
}

func (c *client) receiveMessage() (*sqs.ReceiveMessageOutput, error) {
	param, err := c.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("/lambda/prod/ws-colly/ws-colly-lambda-sqs-url"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	queueURL := *param.Parameter.Value

	resp, err := c.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &queueURL,
		AttributeNames: []*string{
			aws.String("All"),
		},
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *client) uploadToDynamo(resp *sqs.ReceiveMessageOutput) error {
	for i, message := range resp.Messages {
		log.Printf("%v: Message ID: %v\n", i+1, *message.MessageId)
	}

	for _, msg := range resp.Messages {
		item := map[string]*dynamodb.AttributeValue{
			"ID": {
				S: msg.MessageId,
			},
			"Body": {
				S: msg.Body,
			},
			"TimeStamp": {
				S: msg.Attributes["SentTimestamp"],
			},
		}

		// Put item in DynamoDB
		input := &dynamodb.PutItemInput{
			Item:         item,
			TableName:    aws.String("ws-colly-dynamo-table"),
			ReturnValues: aws.String("ALL_OLD"),
		}

		r, err := c.dynamodb.PutItem(input)
		if err != nil {
			log.Printf("Got error calling PutItem: %v", err)
			return err
		}
		log.Printf("PutItem Return: %v\n", r)
	}
	return nil
}

func HandleRequest() (string, error) {
	log.Println("Starting Lambda...")
	c := initClient()
	log.Printf("Client: %v\n", c)

	msgs, err := c.receiveMessage()
	if err != nil {
		log.Printf("Got error calling ReceiveMessage: %v", err)
		return "", err
	}

	log.Println("Messages: ", msgs)

	if err := c.uploadToDynamo(msgs); err != nil {
		log.Printf("Got error calling uploadToDynamo: %v", err)
		return "", err
	}
	log.Printf("Error: %v\n", err)
	log.Printf("Finished Lambda...")

	return "Is this working?", nil
}

func main() {
	lambda.Start(HandleRequest)
}
