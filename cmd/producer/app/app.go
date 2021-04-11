package app

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/streadway/amqp"
	"gorabbit/pkg"
	"log"
	"os"
	"os/signal"
	"time"
)

// NewApp returns App instance
func NewApp(opts ...Option) *App {
	app := &App{
		config: &Config{},
	}

	// Apply custom settings
	for _, opt := range opts {
		opt(app)
	}

	return app
}

// App is just abstraction to keep main file clean and allow testing
type App struct {
	config *Config
}

// Run starts the app, and finished when producer finishes sending all the messages from data provider
func (r *App) Run() error {
	options := &pkg.RabbitProducerOptions{
		ConnectionString: r.config.RabbitMQConnectionString,
		Topic:            r.config.topic,
		ConnectionName:   r.config.connectionName,
		ExchangeType:     r.config.exchangeType,
		UseHeader:        r.config.useHeader,
		MaxReconnects:    r.config.MaxReconnects,
		ReconnectDelay:   r.config.ReconnectDelay,
		MaxRetries:       r.config.MaxSendMessageRetries,
	}

	producer := pkg.NewProducer(options)
	err := producer.Connect()
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
		return err
	}

	producer.InitListeners()

	/*
		If message cannot be routed, RabbitMQ will send it back to you,
		so you can decide what to do.

		For example, implement fallback by storing messages in db, s3 and etc.
	 */
	_ = producer.OnFailedDelivery(func(message amqp.Return) {
		log.Printf("Message was not delivered. Missing bindings? (%s)", message.MessageId)
	})

	/*
		This listener will notify about closed connection.
		The err will be empty if it was triggered by us.
		If you close connection from RabbitMQ UI, or there was any other issue,
		err will include full details.

		We could trigger here reconnect, but in thi POC reconnect is done
		when we are failing to send the messages.
	 */
	_ = producer.OnConnectionError(func(err *amqp.Error) {
		if err != nil {
			log.Printf("Connection closed with error: (%s)", err)
			return
		}

		log.Printf("Connection closed gracefully")
	})

	/*
		When connection closed, all internal amqp lib channels becomes closed as well
		so we need to reinitialize our listeners
	 */
	_ = producer.OnReconnect(func(attempts int) {
		log.Printf("Reconnected on attempt: (%d)", attempts)
		producer.InitListeners()
	})

	/*
		Context used for graceful shutdown.
		We pass it down all the way, so we can react on context change inside SendMessage method
	 */
	ctx, cancel := context.WithCancel(context.Background())
	r.setupGracefulShutdown(cancel, producer)

	/*
		Stream our test file data to RabbitMQ
	 */
	dataProvider := NewDataProvider(r.config.dataSetPath, r.config.dataSetSkipHeader)
	err = dataProvider.Open()
	if err != nil {
		return err
	}

	for row := range dataProvider.Iter() {
		if dataProvider.HasError() {
			panic(dataProvider.Error())
		}

		envelope, err := r.createOrderMessage(row)
		if err != nil {
			log.Printf("Failed to create message: %s", err)
		}

		_, err = producer.SendMessage(ctx, envelope)
		if err != nil {
			log.Printf("App failed to send message to RabbitMQ: %s", envelope.MessageId)
			return err
		}

		log.Printf("Message sent: %s", envelope.MessageId)
		r.delay()
	}

	return nil
}

func (r App) GetTime(str string) *time.Time {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, str)
	if err != nil {
		return nil
	}

	return &t
}

// createOrderMessage maps csv row parsed data to struct
func (r *App) createOrderMessage(row []string) (*pkg.RabbitEnvelope, error) {
	message := &OrderMessage{
		OrderId:     row[0],
		CustomerId:  row[1],
		MessageType: row[2],
	}

	message.PurchaseTimestamp = r.GetTime(row[3])
	message.ApprovedAt = r.GetTime(row[4])
	message.DeliveredCarrierDate = r.GetTime(row[5])
	message.DeliveredCustomerDate = r.GetTime(row[6])
	message.EstimatedDeliveryDate = r.GetTime(row[7])

	envelope := &pkg.RabbitEnvelope{
		RoutingKey: fmt.Sprintf("%s", message.OrderId),
		HashKey:    fmt.Sprintf("%s", message.OrderId),
		Message:    message,
		Immediate:  false,
	}

	id, err := uuid.NewV4()
	if err != nil {
		return envelope, err
	}

	envelope.MessageId = id.String()
	return envelope, nil
}

func (r *App) delay() {
	if r.config.Delay > 0 {
		time.Sleep(time.Duration(r.config.Delay) * time.Millisecond)
	}
}

// setupGracefulShutdown setup listeners for kill signal
// and initiate graceful shutdown
func (r *App) setupGracefulShutdown(cancel context.CancelFunc, producer *pkg.RabbitProducer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Printf("Graceful shutdown triggered by signal: %s", sig)
			cancel()
			producer.TearDown()
			signal.Stop(c)
		}
	}()
}
