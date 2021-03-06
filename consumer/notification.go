// +build notification

package consumer

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/orchid/model/dynamo"
	"github.com/theKono/orchid/model/messagejson"
	"github.com/theKono/orchid/sqs"
)

// consumeNotificationMessage is the a Notification consumer.
//
// It deletes the SQS message regardless of its validity. If it is a
// valid message, then it will insert the Notification both into DynamoDB.
var consumeNotificationMessage = func(message *awsSqs.Message) error {
	var (
		newsFeed *messagejson.Notification
		doc      *dynamodb.PutItemInput
		err      error
	)

	defer sqs.DeleteMessage(message)

	if newsFeed, err = messagejson.NewNotification(message.Body); err != nil {
		log.Println("Cannot parse news feed message\n", err)
		return err
	}

	if doc, err = dynamo.NewNotification(newsFeed); err != nil {
		log.Println("Cannot create news feed dynamodb model\n", err)
		return err
	}

	if err = insertIntoDynamoDB(doc); err != nil {
		return err
	}

	return nil
}

// ConsumeNotification is a function that implements MessageConsumer
// interface.
//
// It is derived from consumeNotificationMessage decorated by
// DecorateConsumeFn.
var ConsumeNotification = DecorateConsumeFn(consumeNotificationMessage)
