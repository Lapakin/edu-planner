package queue

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/Lapakin/edu-planner/internal/adapter/broker"
	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
)

func StartMassageConsumer(ctx context.Context, l *logging.Logger, cfg *broker.Config, processEvent func(context.Context, *domain.Massage) error) error {
	awsCfg, err := cfg.LoadAWSConfig(ctx)
	if err != nil {
		l.Errorf("Failed to load AWS config for SQS consumer: %v", err)
		return err
	}

	client := sqs.NewFromConfig(awsCfg)
	urlResult, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(cfg.QueueName),
	})
	if err != nil {
		l.Errorf("Failed to get queue url for consumer: %v", err)
		return err
	}

	sqsURL := *urlResult.QueueUrl
	go func() {
		for {
			output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(sqsURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			})
			if err != nil {
				l.Errorf("Failed to receive messages from SQS: %v", err)
				continue
			}

			for _, message := range output.Messages {
				var event domain.Massage
				if err = json.Unmarshal([]byte(*message.Body), &event); err != nil {
					l.Errorf("Failed to unmarshal event: %v", err)
					continue
				}

				var processErr error
				for retries := range 3 {
					processErr = processEvent(ctx, &event)
					if processErr == nil {
						break
					}
					l.Warnf("processEvent failed, retrying... (%d/3): %v", retries+1, processErr)
					time.Sleep(time.Second * 2)
				}
				if processErr != nil {
					l.Errorf("Failed to process event after 3 retries: %v", processErr)
					continue
				}

				if _, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(sqsURL),
					ReceiptHandle: message.ReceiptHandle,
				}); err != nil {
					l.Errorf("Failed to delete message from SQS: %v", err)
				}
			}
		}
	}()

	return nil
}
