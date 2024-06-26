// Lambda to parse json from SQS queue and upload to dynamoDB
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

type client struct {
	sqs      *sqs.SQS
	dynamodb *dynamodb.DynamoDB
	ssm      *ssm.SSM
}

func initClients() *client {
	c := &client{}

	s := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	c.sqs = sqs.New(s)
	c.dynamodb = dynamodb.New(s)
	c.ssm = ssm.New(s)

	return c
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

func (c *client) uploadToDynamo(msgs *sqs.ReceiveMessageOutput) error {
	for i, message := range msgs.Messages {
		fmt.Printf("%v: Message ID: %v\n", i+1, *message.MessageId)
	}

	for _, msg := range msgs.Messages {
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
			fmt.Printf("Got error calling PutItem: %v", err)
			return err
		}
		fmt.Printf("PutItem Return: %v\n", r)
	}
	return nil
}

func HandleRequest() error {
	c := initClients()
	fmt.Printf("Client: %v\n", c)

	msgs, err := c.receiveMessage()
	if err != nil {
		fmt.Printf("Got error calling ReceiveMessage: %v", err)
		return err
	}

	fmt.Println("Messages: ", msgs)

	if err := c.uploadToDynamo(msgs); err != nil {
		fmt.Printf("Got error calling uploadToDynamo: %v", err)
		return err
	}
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Finished Lambda...")

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
