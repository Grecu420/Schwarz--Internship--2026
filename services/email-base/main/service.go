package main

import (
	"Schwarz--Internship--2026/services/common"
	"Schwarz--Internship--2026/services/common/email"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/Azure/go-amqp"
)

type EmailConfig struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func sendEmail(p email.EmailMessage, c *EmailConfig) error {
	auth := smtp.PlainAuth("", c.from, c.password, c.smtpHost)
	msg := fmt.Appendf(nil, "To: %s\r\nSubject: %s\r\n\r\n%s", p.To, p.Subject, p.Body)

	return smtp.SendMail(c.smtpHost+":"+c.smtpPort, auth, c.from, []string{p.To}, msg)
}

func readEmailConfig() (*EmailConfig, error) {

	bin, err := os.ReadFile("/run/secrets/email-password")
	if err != nil {
		return nil, fmt.Errorf("failed to read db password secret: %w", err)
	}
	password := strings.TrimSpace(string(bin))

	from, err := common.GetRequiredEnv("FROM")
	if err != nil {
		return nil, err
	}
	port, err := common.GetRequiredEnv("SMTP_PORT")
	if err != nil {
		return nil, err
	}
	host, err := common.GetRequiredEnv("SMTP_HOST")
	if err != nil {
		return nil, err
	}

	return &EmailConfig{password: password, smtpHost: host, smtpPort: port, from: from}, nil

}

func main() {

	rabbitmq_url, err := common.GetRequiredEnv("RABBITMQ_URL")
	if err != nil {
		log.Fatalf("Failed to register url: %v", err)

	}

	emailConf, err := readEmailConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	ctx := context.Background()

	// 1. Connect to RabbitMQ using AMQP 1.0
	conn, err := amqp.Dial(ctx, rabbitmq_url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}
	defer conn.Close()

	// 2. Open an AMQP 1.0 Session
	session, err := conn.NewSession(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}

	// 3. Create a Receiver attached to the target queue
	// Note: The queue "email_queue" must already exist on the broker.
	receiver, err := session.NewReceiver(ctx, "email_queue", nil)
	if err != nil {
		log.Fatalf("Failed to create receiver: %v", err)
	}

	log.Println("Waiting for messages")

	// 4. Consume messages in a loop
	for {
		// Receive blocks until a message is available
		msg, err := receiver.Receive(ctx, nil)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			time.Sleep(1 * time.Second) // Prevent tight loop if connection drops
			continue
		}

		var payload email.EmailMessage
		// GetData() returns the primary payload of the AMQP 1.0 message
		if err := json.Unmarshal(msg.GetData(), &payload); err != nil {
			log.Printf("Malformed message payload: %v", err)

			// Reject message (drops it or sends to Dead Letter Queue if configured)
			// The third argument is an optional *amqp.Error mapping.
			_ = receiver.RejectMessage(ctx, msg, nil)
			continue
		}

		if err := sendEmail(payload, emailConf); err != nil {
			log.Printf("Email send failed to %s: %v", payload.To, err)

			// Release message (equivalent to Nack with requeue=true)
			// so it can be picked up by another worker or retried.
			_ = receiver.ReleaseMessage(ctx, msg)
			continue
		}

		log.Printf("Email successfully sent to %s", payload.To)

		// Accept the message (equivalent to Ack) to mark it as completed
		if err := receiver.AcceptMessage(ctx, msg); err != nil {
			log.Printf("Failed to accept message: %v", err)
		}
	}
}
