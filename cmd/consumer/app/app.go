package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"gorabbit/pkg"
	"log"
	"os"
	"os/signal"
	"time"
)

// NewApp returns new instance of the Consumer App
func NewApp(opts ...Option) *App {
	app := &App{
		config: &Config{},
	}

	for _, opt := range opts {
		opt(app)
	}

	// just auto generate names for connections and consumer tag to be used for cancellation on graceful shut down
	app.config.InstanceId = fmt.Sprintf("%s_%d", app.config.Identity, os.Getpid())
	return app
}

// Config for app
type Config struct {
	RabbitMQConnectionString string   // Connection string to RabbitMQ
	QueueNames               []string // queues that consumer will subscribe to
	Identity                 string   // for metadata in logs and other tools
	InstanceId               string   // id for identify connection in RabbitMQ UI
	SingleActiveConsumer     bool     // RabbitMQ will ensure only a single consumer receives the messages
	PrefetchSize             int      // Number of messages pushed by RabbitMQ to the consumer and kept in buffer
	BindingKey               string   // used for binding queues to exchange
	Topic                    string
	UseHeader                bool
}

// App represents Consumer App to keep main file clean
type App struct {
	config *Config
}

// Run starts the app
func (r *App) Run() error {
	// Create a root context
	ctx := context.Background()
	// Wrap root context with a new context, with its cancellation function
	// from the original context
	ctx, cancel := context.WithCancel(ctx)

	options := &pkg.RabbitConsumerOptions{
		ConnectionString:     r.config.RabbitMQConnectionString,
		QueueNames:           r.config.QueueNames,
		ConnectionName:       r.config.InstanceId,
		Identity:       r.config.Identity,
		Topic:                r.config.Topic,
		SingleActiveConsumer: r.config.SingleActiveConsumer,
		PrefetchSize:         r.config.PrefetchSize,
		InstanceId:           r.config.InstanceId,
		BindingKey:           r.config.BindingKey,
		UseHeader:           r.config.UseHeader,
	}

	consumer := pkg.NewConsumer(options)

	err := consumer.Connect()
	if err != nil {
		log.Printf("Failed to connect: %s", err)
		return err
	}

	err = consumer.InitQueues()
	if err != nil {
		log.Printf("Failed to init queues: %s", err)
		return err
	}

	consumer.InitListeners()

	// We want gracefully shut down the consumer,
	// to allow us finish messages processing or safely stop at some point,
	// so we setup listeners when to trigger graceful shutdown
	r.setupGracefulShutdown(cancel, consumer)

	err = consumer.Subscribe(ctx, func(ctx context.Context, message amqp.Delivery) error {
		queueName, _ := ctx.Value("queueName").(string)
		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			log.Printf("Oh, I see someone asked for shut down, I'll not event start the work")
			return err
		}

		log.Printf("%s - %s", queueName, string(message.Body))
		for i := 0; i < 10; i++ {
			log.Printf("I'm still processing the message: %d -> %s -> %s", i, queueName, message.Body)
			time.Sleep(1000 * time.Millisecond)

			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				log.Printf("Oh, I see someone asked for shut down, I'm in the middle of heavy procssing, but I know to wtop safely here...")
				return err
			}
		}

		return nil
	})

	return err
}

// setupGracefulShutdown setup listeners for kill signals or connection error
// and initiate graceful shutdown
func (r *App) setupGracefulShutdown(cancel context.CancelFunc, consumer *pkg.RabbitConsumer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Printf("Graceful shutdown triggered by signal: %s", sig)
			cancel()
			consumer.TearDown()
			signal.Stop(c)
		}
	}()

	go func() {
		for err := range consumer.ConnectionError {
			log.Printf("Graceful shutdown triggered by network error: %s", err)
			cancel()
			consumer.TearDown()
		}
	}()
}
