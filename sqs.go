package sqs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/sqs", new(SQS))
}

type SQS struct{}

type Client interface {
	QueueURL(queueName string) string
	CreateQueue(queueName string) (*sqs.CreateQueueOutput, error)
	DeleteQueue(queueName string) (*sqs.DeleteQueueOutput, error)
	SendMessage(kw SendMessageModel) (*sqs.SendMessageOutput, error)
	ReceiveMessage(kw ReceiveMessageModel) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(kw DeleteMessageModel) (*sqs.DeleteMessageOutput, error)
	SendMessageBatch(kw SendMessageBatchModel) (*sqs.SendMessageBatchOutput, error)
	DeleteMessageBatch(kw DeleteMessageBatchModel) (*sqs.DeleteMessageBatchOutput, error)
}
type client struct {
	session *sqs.SQS
	UserID  string
	Url     string
}

var instance *client

func (s *SQS) New(opts *ClientModel) (Client, error) {
	if instance == nil {
		if opts.Url == "" || opts.UserID == "" {
			return nil, errors.New("Parametrs `Url` and `UserID` is required")
		}
		if opts.AccessKeyID == "" {
			opts.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		}
		if opts.SecretAccessKey == "" {
			opts.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		}
		if opts.Region == "" {
			opts.Region = os.Getenv("AWS_DEFAULT_REGION")
		}

		if opts.AccessKeyID == "" || opts.SecretAccessKey == "" || opts.Region == "" {
			log.Fatal("AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY or AWS_DEFAULT_REGION not found in environment")
		}

		httpClient, err := NewHTTPClientWithSettings(HTTPClientSettings{
			Connect:          5 * time.Second,
			ExpectContinue:   1 * time.Second,
			IdleConn:         90 * time.Second,
			ConnKeepAlive:    0 * time.Second,
			MaxAllIdleConns:  200,
			MaxHostIdleConns: 20,
			ResponseHeader:   5 * time.Second,
		})

		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(opts.Region),
			Credentials: credentials.NewStaticCredentials(opts.AccessKeyID, opts.SecretAccessKey, ""),
			Endpoint:    aws.String(opts.Url),
			HTTPClient:  httpClient,
		})

		if err != nil {
			log.Fatal(err)
		}
		instance = new(client)
		instance.session = sqs.New(sess)
		instance.UserID = opts.UserID
		instance.Url = opts.Url
	}
	return instance, nil
}

func (c *client) QueueURL(queueName string) string {
	return fmt.Sprintf("%s/%s/%s", c.Url, c.UserID, queueName)
}

func (c *client) CreateQueue(queueName string) (*sqs.CreateQueueOutput, error) {
	params := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}

	resp, err := c.session.CreateQueue(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) DeleteQueue(queueName string) (*sqs.DeleteQueueOutput, error) {
	params := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(c.QueueURL(queueName)),
	}

	resp, err := c.session.DeleteQueue(params)
	if err != nil {
		return nil, err
	}
	return resp, err
}

//Delivers a message to the specified queue
func (c *client) SendMessage(kw SendMessageModel) (*sqs.SendMessageOutput, error) {
	params := &sqs.SendMessageInput{
		QueueUrl:          aws.String(c.QueueURL(*kw.QueueName)),
		MessageBody:       aws.String(*kw.MessageBody),
		MessageAttributes: kw.MessageAttributes,
	}
	if kw.DelaySeconds != nil {
		params.DelaySeconds = aws.Int64(*kw.DelaySeconds)
	}

	resp, err := c.session.SendMessage(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Retrieves one or more messages, with a maximum limit of 10 messages, from
// the specified queue.
func (c *client) ReceiveMessage(kw ReceiveMessageModel) (*sqs.ReceiveMessageOutput, error) {
	queueUrl := c.QueueURL(*kw.QueueName)
	params := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueUrl),
		AttributeNames:        kw.AttributeNames,
		MessageAttributeNames: kw.MessageAttributeNames,
	}
	if kw.MaxNumberOfMessages != nil {
		params.MaxNumberOfMessages = aws.Int64(*kw.MaxNumberOfMessages)
	}
	if kw.WaitTimeSeconds != nil {
		params.WaitTimeSeconds = aws.Int64(*kw.WaitTimeSeconds)
	}
	if kw.VisibilityTimeout != nil {
		params.VisibilityTimeout = aws.Int64(*kw.VisibilityTimeout)
	}

	resp, err := c.session.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Deletes the specified message from the specified queue.
func (c *client) DeleteMessage(kw DeleteMessageModel) (*sqs.DeleteMessageOutput, error) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.QueueURL(*kw.QueueName)),
		ReceiptHandle: aws.String(*kw.ReceiptHandle),
	}

	resp, err := c.session.DeleteMessage(params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Delivers up to ten messages to the specified queue.
func (c *client) SendMessageBatch(kw SendMessageBatchModel) (*sqs.SendMessageBatchOutput, error) {
	params := &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(c.QueueURL(*kw.QueueName)),
		Entries:  kw.Entries,
	}

	resp, err := c.session.SendMessageBatch(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//Deletes up to ten messages from the specified queue.
func (c *client) DeleteMessageBatch(kw DeleteMessageBatchModel) (*sqs.DeleteMessageBatchOutput, error) {
	params := &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(c.QueueURL(*kw.QueueName)),
		Entries:  kw.Entries,
	}

	resp, err := c.session.DeleteMessageBatch(params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
