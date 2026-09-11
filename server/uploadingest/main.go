package uploadingest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	uuid "github.com/google/uuid"
)

type VideoStore interface {
	MarkVideoReady(publicID uuid.UUID) error
}

type s3Event struct {
	Records []struct {
		EventName string `json:"eventName"`
		S3        struct {
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func Run(ctx context.Context, store VideoStore, queue string, endpoint string) error {
	queueURL := endpoint + "/" + queue
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"key", "secret", "sessionstring",
		)),
	)
	if err != nil {
		return err
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	slog.Info("starting upload event consumer", "queue", queueURL)
	for {
		select {
		case <-ctx.Done():
			slog.Info("upload event consumer stopped")
			return nil
		default:
		}

		out, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(queueURL),
			MaxNumberOfMessages:   10,
			WaitTimeSeconds:       20,
			VisibilityTimeout:     30,
			MessageAttributeNames: []string{"All"},
		})
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.Error("failed to receive messages", "error", err)
			// Retry slowly to avoid flooding logs on startup
			timer := time.NewTimer(5 * time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil
			case <-timer.C:
			}
			continue
		}

		for _, msg := range out.Messages {
			slog.Info("received upload queue message", "message_id", aws.ToString(msg.MessageId))
			if handle(ctx, store, msg) {
				if _, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueURL),
					ReceiptHandle: msg.ReceiptHandle,
				}); err != nil {
					slog.Error("failed to delete message", "error", err, "message_id", aws.ToString(msg.MessageId))
				} else {
					slog.Info("deleted upload queue message", "message_id", aws.ToString(msg.MessageId))
				}
			}
		}
	}
}

func handle(ctx context.Context, store VideoStore, msg types.Message) bool {
	var event s3Event
	if err := json.Unmarshal([]byte(aws.ToString(msg.Body)), &event); err != nil {
		slog.Error("malformed message body, dropping", "error", err, "message_id", aws.ToString(msg.MessageId))
		return true
	}

	if len(event.Records) == 0 {
		slog.Info("ignoring non-upload event", "body", aws.ToString(msg.Body))
		return true
	}

	ok := true
	for _, record := range event.Records {
		publicID, err := videoIDFromKey(record.S3.Object.Key)
		if err != nil {
			slog.Error("unrecognized object key, skipping", "key", record.S3.Object.Key, "error", err)
			continue
		}

		if err := store.MarkVideoReady(publicID); err != nil {
			slog.Error("failed to mark video ready", "error", err, "public_id", publicID)
			ok = false
		} else {
			slog.Info("upload event processed", "message_id", aws.ToString(msg.MessageId), "public_id", publicID, "key", record.S3.Object.Key)
		}
	}
	return ok
}

func videoIDFromKey(key string) (uuid.UUID, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 3 || parts[0] != "videos" || !strings.HasSuffix(parts[2], ".mp4") {
		return uuid.Nil, errors.New("unexpected s3 key layout")
	}
	return uuid.Parse(parts[1])
}
