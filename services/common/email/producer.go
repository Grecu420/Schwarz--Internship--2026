package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/go-amqp"
)

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Producer struct {
	conn    *amqp.Conn
	session *amqp.Session
	sender  *amqp.Sender
}

// NewProducer initializes a long-lived AMQP 1.0 connection and sender.
func NewProducer(ctx context.Context, amqpURL, queueName string) (*Producer, error) {
	conn, err := amqp.Dial(ctx, amqpURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AMQP broker: %w", err)
	}

	session, err := conn.NewSession(ctx, nil)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open AMQP session: %w", err)
	}

	sender, err := session.NewSender(ctx, queueName, nil)
	if err != nil {
		_ = session.Close(ctx)
		_ = conn.Close()
		return nil, fmt.Errorf("failed to create AMQP sender: %w", err)
	}

	return &Producer{
		conn:    conn,
		session: session,
		sender:  sender,
	}, nil
}

// Publish serializes and sends an email payload to the queue. Safe for concurrent use.
func (p *Producer) Publish(ctx context.Context, payload EmailMessage) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	msg := amqp.NewMessage(data)
	msg.Header = &amqp.MessageHeader{
		Durable: true, // Persists message on disk across broker restarts
	}

	if err := p.sender.Send(ctx, msg, nil); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// Close gracefully closes the sender, session, and connection.
func (p *Producer) Close(ctx context.Context) error {
	_ = p.sender.Close(ctx)
	_ = p.session.Close(ctx)
	return p.conn.Close()
}
