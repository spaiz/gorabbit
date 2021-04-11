package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

/*
	NewProducer returns new instance of RabbitProducer instance
 */
func NewProducer(options *RabbitProducerOptions) *RabbitProducer {
	producer := &RabbitProducer{
		options: options,
		reconnectGuard:   &sync.Mutex{},
		reconnected:      make(chan int),
		connectionError:  make(chan *amqp.Error),
		undeliverable:    make(chan amqp.Return),
	}

	return producer
}

type RabbitProducerOptions struct {
	ConnectionString string
	Name             string
	Topic            string
	ConnectionName   string
	ExchangeType     string
	UseHeader        bool
	MaxReconnects    int
	ReconnectDelay   time.Duration
	MaxRetries       int
}

type RabbitProducer struct {
	isStopped            bool
	isTearDownInProgress bool
	conn                 *amqp.Connection
	ch                   *amqp.Channel
	isReconnecting       bool
	reconnectGuard       *sync.Mutex
	undeliverable        chan amqp.Return
	mandatory            bool
	tearDownOnce         sync.Once
	connectionError      chan *amqp.Error
	reconnected          chan int
	options              *RabbitProducerOptions
}

/*
	Connect opens connection to RabbitMQ
	There is no retries. We must be able to connect at least once.
 */
func (r *RabbitProducer) Connect() error {
	conn, err := amqp.DialConfig(r.options.ConnectionString, amqp.Config{
		Properties: amqp.Table{"connection_name": r.options.ConnectionName,},
	})

	if err != nil {
		log.Printf("Failed to open connection to RabbitMQ: %s", err)
		return err
	}

	r.conn = conn

	r.ch, err = conn.Channel()
	if err != nil {
		log.Printf("Failed to create channel: %d", err)
		return err
	}

	/*
		Declare exchange. It's durable, so will survive RabbitMQ restarts.
		It will fail if the exchange is already exists with different settings, otherwise succeed
	 */
	err = r.ch.ExchangeDeclare(r.options.Topic, r.options.ExchangeType, true, false, false, false, amqp.Table{})
	if err != nil {
		return err
	}

	return err
}

/*
	InitListeners set ups listeners for connection (closed connection) and delivery errors
 */
func (r *RabbitProducer) InitListeners() {
	undeliverable := r.ch.NotifyReturn(make(chan amqp.Return))

	go func() {
		for str := range undeliverable {
			r.undeliverable <- str
		}
	}()

	connectionError := r.ch.NotifyClose(make(chan *amqp.Error))

	go func() {
		for err := range connectionError {
			r.connectionError <- err
		}
	}()

	log.Printf("Listeners initialized")
}

// OnFailedDelivery gives you ability to catch situations when no binding exists
// and messages will be just dropped by RabbitMQ
// When callback provided, it sends messages with mandatory flag on: https://www.rabbitmq.com/publishers.html
func (r *RabbitProducer) OnFailedDelivery(callback func(message amqp.Return)) error {
	r.mandatory = true

	go func() {
		for message := range r.undeliverable {
			callback(message)
		}
	}()

	return nil
}

// OnConnectionError will call the callback function when connection error (closed connection) occurred
func (r *RabbitProducer) OnConnectionError(callback func(err *amqp.Error)) error {
	go func() {
		for err := range r.connectionError {
			callback(err)
		}
	}()

	return nil
}

// OnReconnect call callback when reconnect succeed
func (r *RabbitProducer) OnReconnect(callback func(attempts int)) error {
	go func() {
		for err := range r.reconnected {
			callback(err)
		}
	}()

	return nil
}

/*
	SendMessage is doing some validations to be able to catch calling it when tear shut down initialized.
	Connections will be closed, and we don't want to reconnect if we want shut down the instance.
 */
func (r *RabbitProducer) SendMessage(ctx context.Context, envelope *RabbitEnvelope) (*SendMessageResult, error) {
	res := &SendMessageResult{}
	res.ExchangeName = r.options.Topic
	res.Mandatory = r.mandatory
	res.MessageId = envelope.MessageId
	res.Immediate = envelope.Immediate

	// early return if tear down initialized
	if err := r.allowSendMessage(); err != nil {
		return res, err
	}

	body, err := envelope.GetBody()
	if err != nil {
		return res, err
	}

	// We use envelope.HashKey inside the header, cause we want RqabitMQ consistent hash exchange
	// will apply hash function and route the message to relevant queue
	headers := amqp.Table{}
	if r.options.UseHeader {
		headers["hash-on"] = envelope.HashKey
	}

	pb := &amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     headers,
		MessageId:   envelope.MessageId,
	}

	// this is hack to implement dropping messages if there is no consumers on the queues
	// we use TTL because RabbitMQ abandoned support for native immediate function protocol implementation
	// because of complexity
	if envelope.Immediate {
		pb.Expiration = "0"
	}

	// We must protect this method to prevent multiple reconnect attempts from different goroutines
	r.reconnectGuard.Lock()
	defer r.reconnectGuard.Unlock()

	err = r.publish(ctx, envelope, pb)
	return res, err
}

func (r *RabbitProducer) allowSendMessage() error {
	if r.isTearDownInProgress {
		return errors.New("tear down already in progress")
	}

	if r.isStopped {
		return errors.New("message will not be sent, producer stopped")
	}

	return nil
}

/*
	publish method actually tries to publish the message to RabbitMQ
	It can be configured to do retries on ErrClosed errors,
	and will early return on context cancellation or/and tear down triggered
 */
func (r *RabbitProducer) publish(ctx context.Context, envelope *RabbitEnvelope, pb *amqp.Publishing, ) error {
	attempts := 0

	for {
		attempts++
		if err := r.allowSendMessage(); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := r.ch.Publish(
			r.options.Topic,
			envelope.RoutingKey, // routing key
			r.mandatory,
			false,
			*pb,
		)

		if err == nil {
			return nil
		}

		if err != amqp.ErrClosed {
			return nil
		}

		if attempts > r.options.MaxRetries {
			log.Printf("SendMessage retries exeeced: %d > %d", attempts, r.options.MaxRetries)
			return err
		}

		log.Printf("Failed to send message on attempt [%d]: %s", attempts, err)

		r.isReconnecting = true
		err = r.reconnect(ctx)
		r.isReconnecting = false

		if err != amqp.ErrClosed {
			return err
		}
	}

	return nil
}

/*
	reconnect method tries to reconnect RabbitMQ server.
	It has own retries options and early return on context cancellation or/and tear down trigger
 */
func (r *RabbitProducer) reconnect(ctx context.Context) error {
	attempt := 0

	var err error
	for attempt < r.options.MaxReconnects {
		attempt += 1
		log.Printf("Reconnecting attempt: %d", attempt)
		err = r.Connect()
		if err == nil {
			select {
			case r.reconnected <- attempt:
			default:
				log.Printf("No listener for [reconnected] event")
			}

			log.Printf("Succesfully reconnected")
			return nil
		}

		if err != amqp.ErrClosed {
			return err
		}

		if err = r.allowSendMessage(); err != nil {
			return err
		}


		// Some delay before reconnecting
		delay := time.NewTimer(r.options.ReconnectDelay)

		// If context was cancelled we want early return and do not wait for delay to finish
		select {
		case <-delay.C:
			log.Printf("Delay ended, next try")
			delay.Stop()
		case <-ctx.Done():
			delay.Stop()
			log.Printf("Reconnect cancelled, return")
			return ctx.Err()
		default:

		}
	}

	return err
}

// TearDown deals with graceful shutdown
// The driver will ensure we finish sending the messages already in progress to not lose them
// And than safely close the connection
// If someone will try send message on stopped producer, we will return error
func (r *RabbitProducer) TearDown() error {
	r.tearDownOnce.Do(func() {
		r.isTearDownInProgress = true
		defer func() {
			r.isTearDownInProgress = false
			r.isStopped = true
		}()

		log.Printf("close the connection")
		if err := r.conn.Close(); err != nil {
			log.Printf("attempt to close RabbitMQ connection failed: %s", err)
			return
		}

		defer log.Printf("Producer teardown succesfully completed")
	})

	return nil
}

// RabbitEnvelope wraps our message with some rabbit specific options
//
// RoutingKey - routing key used by exchange to decide were message should be routed
// Immediate - when true, will set the TTL to 0
// and cause messages to be expired upon reaching a queue
// unless they can be delivered to a consumer immediately.
// Done with TTL instead of native Immediate flag because of this - http://www.rabbitmq.com/blog/2012/11/19/breaking-things-with-rabbitmq-3-0/
// Message - Message struct that will be marshalled to JSON and represents actual payload
// HashKey - value to be used by consistent hash exchange
// MessageId - message id
// Mandatory - This flag tells the server how to react if a message cannot be routed to a queue.
// Specifically, if mandatory is set and after running the bindings the message was placed on zero
//queues then the message is returned to the sender (with a basic.return).
//If mandatory had not been set under the same circumstances the server would silently drop the message.
type RabbitEnvelope struct {
	RoutingKey string
	Message    interface{}
	Immediate  bool // when true, message will dropped when no consumers on the queue
	HashKey    string
	MessageId  string
}

// GetBody returns json representation of our message
func (r *RabbitEnvelope) GetBody() ([]byte, error) {
	return json.Marshal(r.Message)
}

// SendMessageResult used as DTO for sending message results
type SendMessageResult struct {
	ExchangeName string
	Mandatory    bool
	MessageId    string
	Immediate    bool
}
