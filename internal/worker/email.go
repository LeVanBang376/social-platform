package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"social-platform/internal/service/email"

	"github.com/redis/go-redis/v9"
)

const (
	emailStream = "email_jobs"
	emailGroup  = "email-workers"
)

type EmailWorker struct {
	redisClient  *redis.Client
	emailService *email.EmailService
}

func NewEmailWorker(
	redisClient *redis.Client,
	emailService *email.EmailService,
) *EmailWorker {
	return &EmailWorker{
		redisClient:  redisClient,
		emailService: emailService,
	}
}

func (w *EmailWorker) Start(
	ctx context.Context,
	consumerName string,
) {
	// Create consumer group if it does not exist.
	err := w.redisClient.XGroupCreateMkStream(
		ctx,
		emailStream,
		emailGroup,
		"$",
	).Err()

	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Printf(
			"failed to create email consumer group: %v",
			err,
		)
		return
	}

	log.Printf(
		"email worker started: %s",
		consumerName,
	)

	for {
		messages, err := w.redisClient.XReadGroup(
			ctx,
			&redis.XReadGroupArgs{
				Group:    emailGroup,
				Consumer: consumerName,
				Streams:  []string{emailStream, ">"},
				Count:    1,
				Block:    0,
			},
		).Result()

		if err != nil {
			if ctx.Err() != nil {
				log.Printf(
					"email worker stopped: %s",
					consumerName,
				)
				return
			}

			log.Printf(
				"failed to read email job: %v",
				err,
			)

			continue
		}

		for _, stream := range messages {
			for _, message := range stream.Messages {
				if err := w.process(
					ctx,
					message,
				); err != nil {
					log.Printf(
						"failed to process email job %s: %v",
						message.ID,
						err,
					)

					continue
				}

				if err := w.redisClient.XAck(
					ctx,
					emailStream,
					emailGroup,
					message.ID,
				).Err(); err != nil {
					log.Printf(
						"failed to ACK email job %s: %v",
						message.ID,
						err,
					)

					continue
				}

				log.Printf(
					"email job completed: %s",
					message.ID,
				)
			}
		}
	}
}

func (w *EmailWorker) process(
	ctx context.Context,
	message redis.XMessage,
) error {
	rawData, ok := message.Values["data"].(string)
	if !ok {
		return fmt.Errorf("invalid email job data")
	}

	var mes email.EmailMessage

	if err := json.Unmarshal(
		[]byte(rawData),
		&mes,
	); err != nil {
		return err
	}

	if err := w.emailService.Send(
		mes.To,
		mes.Subject,
		mes.Body,
	); err != nil {
		return err
	}

	return nil
}
