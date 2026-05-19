package queue

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"

	adapterJson "github.com/Lapakin/edu-planner/internal/adapter/json"
)

func StartMassageProducers(queueName string, dbURL string, db *sqlx.DB, l *logging.Logger, cfg aws.Config) error {
	client := sqs.NewFromConfig(cfg)
	urlResult, err := client.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		l.Errorf("Failed to get queue url: %v", err)
		return err
	}

	sqsURL := *urlResult.QueueUrl
	listener := pq.NewListener(dbURL, 10*time.Second, time.Minute, func(_ pq.ListenerEventType, err error) {
		if err != nil {
			l.Errorf("Listener error: %v", err)
		}
	})

	if err = listener.Listen("events_channel"); err != nil {
		l.Errorf("Failed to listen to events_channel: %v", err)
		return err
	}

	l.Infoln("Started event producer listening to PostgreSQL...")
	go func() {
		for {
			select {
			case n := <-listener.Notify:
				if n == nil {
					continue
				}

				var msg domain.Massage
				var objectBytes []byte
				row := db.QueryRowContext(context.Background(),
					`SELECT event_id, action_type, object_type, object
					 FROM massage WHERE event_id = $1`, n.Extra)
				if scanErr := row.Scan(&msg.EventID, &msg.ActionType, &msg.ObjectType, &objectBytes); scanErr != nil {
					l.Errorf("Failed to fetch massage row (event_id=%s): %v", n.Extra, scanErr)
					continue
				}
				msg.Object = adapterJson.RawMessage(objectBytes)

				body, marshalErr := adapterJson.Marshal(msg)
				if marshalErr != nil {
					l.Errorf("Failed to marshal massage: %v", marshalErr)
					continue
				}

				_, err = client.SendMessage(context.Background(), &sqs.SendMessageInput{
					MessageBody: aws.String(string(body)),
					QueueUrl:    aws.String(sqsURL),
				})
				if err != nil {
					l.Errorf("Failed to send message to SQS: %v", err)
				} else {
					l.Debugln("Message sent to SQS successfully")
				}
			case <-time.After(90 * time.Second):
				go func() { _ = listener.Ping() }()
			}
		}
	}()

	return nil
}
