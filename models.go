package sqs

import "github.com/aws/aws-sdk-go/service/sqs"

type ClientModel struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	UserID          string
	Url             string
}

type SendMessageModel struct {
	QueueName    *string
	MessageBody  *string
	DelaySeconds *int64
	//Each message attribute consists of a Name, Type, and Value.
	MessageAttributes map[string]*sqs.MessageAttributeValue
}

type ReceiveMessageModel struct {
	QueueName           *string
	MaxNumberOfMessages *int64
	WaitTimeSeconds     *int64
	VisibilityTimeout   *int64
	//A list of attributes that need to be returned along with each message
	//e.g. ApproximateNumberOfMessages, ApproximateNumberOfMessagesDelayed
	AttributeNames []*string
	//The name of the message attribute, where N is the index. The message attribute
	MessageAttributeNames []*string
}

type DeleteMessageModel struct {
	QueueName     *string
	ReceiptHandle *string
}

type SendMessageBatchModel struct {
	QueueName *string
	//Contains the details of a single Amazon SQS message along with a Id.
	Entries []*sqs.SendMessageBatchRequestEntry
}

type DeleteMessageBatchModel struct {
	QueueName *string
	Entries   []*sqs.DeleteMessageBatchRequestEntry
}
