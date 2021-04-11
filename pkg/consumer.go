package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

// NewConsumer returns new RabbitConsumer instance
func NewConsumer(options *RabbitConsumerOptions) *RabbitConsumer {
	cons := &RabbitConsumer{
		options:         options,
		done:            make(chan bool),
		cancel:          make(chan string),
		ConnectionError: make(chan *amqp.Error),
	}
	return cons
}

// RabbitConsumerOptions allows setup consumer with different settings
type RabbitConsumerOptions struct {
	ConnectionString     string
	QueueNames           []string
	ConnectionName       string
	Topic                string
	SingleActiveConsumer bool
	PrefetchSize         int
	InstanceId           string
	BindingKey           string
	UseHeader            bool
	Identity             string
}

// MessageSource
type MessageSource struct {
	QueueName   string
	Messages    <-chan amqp.Delivery
	ConsumerTag string
}

// RabbitConsumer abstracts work with RabbitMQ
// and implements all the logic related to life cycle of the consumer
// from creating and binding queues up to graceful shutdown
// Hides all the complex logic form the client that will need only provide options and message handler
type RabbitConsumer struct {
	conn                 *amqp.Connection
	ch                   *amqp.Channel
	messages             <-chan amqp.Delivery
	isTearDownInProgress bool
	isSuccessfully       bool
	sources              []*MessageSource
	handlersDone         sync.WaitGroup
	done                 chan bool
	cancel               chan string
	options              *RabbitConsumerOptions
	tearDownOnce         sync.Once
	ConnectionError      chan *amqp.Error
	UseHeader            bool
}

// Connect creates connection to RabbitMQ
// Creates queues and binds it to specified exchange
// Setups listeners for network and and other RabbitMQ errors
func (r *RabbitConsumer) Connect() error {
	conn, err := amqp.DialConfig(r.options.ConnectionString, amqp.Config{
		Properties: amqp.Table{"connection_name": r.options.ConnectionName,},
	})

	if err != nil {
		log.Printf("Failed to connect RabbitMQ: %s", err)
		return err
	}

	r.conn = conn

	r.ch, err = conn.Channel()
	if err != nil {
		log.Printf("Failed to create channel to RabbitMQ: %s", err)
		return err
	}

	return nil
}

// InitListeners creates listener for connection close errors
func (r *RabbitConsumer) InitListeners() {
	connectionError := r.ch.NotifyClose(make(chan *amqp.Error))

	go func() {
		for err := range connectionError {
			r.ConnectionError <- err
		}
	}()

	log.Printf("Listeners initialized")
}

// InitQueues declares queues and binds it to the exchange
func (r *RabbitConsumer) InitQueues() error {
	/*
		Declare main exchange.
		Think about it like kafka's "Topic"
		It will fail if the exchange is already exists with different settings, otherwise succeed
	*/
	err := r.ch.ExchangeDeclare(r.options.Topic, "fanout", true, false, false, false, amqp.Table{})
	if err != nil {
		log.Printf("Failed to declared exchange: %s", r.options.Topic)
		return err
	}

	/*
		Declare exchange for specific service (identity)
		Messages from the main topic will be copied to this exchange.
		It will fail if the exchange is already exists with different settings, otherwise succeed
	*/

	// Say RabbitMQ we want use header for routing messages with consistent hash exchange type
	args := amqp.Table{}
	if r.options.UseHeader {
		args["hash-header"] = "hash-on"
	}

	exchangeName := fmt.Sprintf("%s_%s", r.options.Topic, r.options.Identity)
	err = r.ch.ExchangeDeclare(exchangeName, "x-consistent-hash", true, false, false, false, args)
	if err != nil {
		log.Printf("Failed to declared exchange: %s", exchangeName)
		return err
	}

	err = r.ch.ExchangeBind(exchangeName, "", r.options.Topic, false, amqp.Table{})
	if err != nil {
		log.Printf("Failed to bind exchanges: %s -> %s", exchangeName, r.options.Topic)
		return err
	}

	queueArgs := amqp.Table{}
	/*
		Single Active Consumer will ensure there is only a single consumer which receives the messages
		This is new cool feature RabbitMQ v3.8 provides: https://www.cloudamqp.com/blog/rabbitmq-3-8-feature-focus-single-active-consumer.html
	 */
	if r.options.SingleActiveConsumer {
		queueArgs["x-single-active-consumer"] = r.options.SingleActiveConsumer
	}


	// Declare queues. This operation is idempotent
	// When declared, attempt to redeclare with different settings will fail
	for _, queueName := range r.options.QueueNames {
		queue, err := r.ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			queueArgs,
		)

		log.Printf("Queue declared: %s, Current amount of consumers: %d, Current amount messages: %d", queue.Name, queue.Consumers, queue.Messages)
		if err != nil {
			log.Printf("Failed to declare queue: %s (%s)", queueName, err)
			return err
		}
	}

	// Bind queues to provided exchange. This operation is idempotent.
	// When bound, attempt to bind with different settings will fail
	for _, queueName := range r.options.QueueNames {
		err := r.ch.QueueBind(queueName, r.options.BindingKey, exchangeName, false, nil)
		if err != nil {
			log.Printf("Failed to bind the queue: [queue]%s -> [key]%s -> [exchange]%s (%s)", queueName, r.options.BindingKey, exchangeName, err)
			return err
		}

		log.Printf("Queue succesfully bound: [queue]%s -> [key]%s -> [exchange]%s", queueName, r.options.BindingKey, exchangeName)
	}

	return nil
}

// Subscribe setups handlers queue
// Currently it creates worker per queue. Although GO's routines are green threads and very lightweight
// it can be overhead when queues are always has messages, and number of workers much more then cores
// Can be optimised if needed by introducing pool of workers
func (r *RabbitConsumer) Subscribe(ctx context.Context, handler func(ctx context.Context, message amqp.Delivery) error) error {
	// Defines how many messages will be pushed from RabbitMQ to consumer
	// Trade off between stability, performance, throughput
	// IMHO most stable is when prefetch size is 1
	err := r.ch.Qos(r.options.PrefetchSize, 0, false)
	if err != nil {
		log.Printf("Failed to set QoS: [PrefetchSize]%d (%s)", r.options.PrefetchSize, err)
		return err
	}

	for i, queueName := range r.options.QueueNames {
		// Identifier of the rabbit's consumer that will be used for ack/nack and other operations
		// Must be unique per opened channel
		consumerTag := fmt.Sprintf("concumer_%s_i%d", r.options.InstanceId, i)
		// autoAck has some memory issues - https://github.com/streadway/amqp/issues/484
		messages, err := r.ch.Consume(queueName, consumerTag, false, false, true, true, nil)
		if err != nil {
			log.Printf("Failed to set up consumer: [queueName]%s -> [consumerTag]%s (%s)", queueName, consumerTag, err)
			return err
		}

		log.Printf("Consumer created: [queue]: %s -> [consumer]: %s", queueName, consumerTag)

		source := &MessageSource{
			QueueName:   queueName,
			Messages:    messages,
			ConsumerTag: consumerTag,
		}

		r.sources = append(r.sources, source)
	}

	// to graceful shutdown, we will wait for waiting group to be done
	// Each worker will responsible to call Done when it finished
	r.handlersDone.Add(len(r.sources))

	for _, source := range r.sources {
		// Starting worker as a simple goroutine without additional abstraction
		log.Printf("Worker started: [queue]: %s -> [consumer]: %s", source.QueueName, source.ConsumerTag)
		go func(source *MessageSource) {
			endedWithCancellation := 0
			ctx := context.WithValue(ctx, "queueName", source.QueueName)
			// wait for message on the opened channel
			// it will automatically exit the loop when channel is closed byt the driver
			// for example as result of graceful shutdown that triggers closing connection
			//
			// Includes minimal logic for handling ack/nack fo the messages depending ont he result of the handler
			for d := range source.Messages {
				err := handler(ctx, d)
				if err == nil {
					err = d.Ack(false)
					if err != nil {
						log.Printf("Cannot ack the message, it's very bad. Ended with concellation? - %d", endedWithCancellation)
						panic(err)
					}

					continue
				}

				if errors.Is(err, context.Canceled) {
					endedWithCancellation++
					log.Printf("Hanlder ended with cancellation: %d", endedWithCancellation)
				} else {
					log.Printf("Hanlder ended with error: %s", err)
				}

				log.Printf("Nack the messages [queue]: %s - [consumer]: %s", source.QueueName, source.ConsumerTag)
				err = d.Nack(false, true)
				if err != nil {
					log.Printf("Cannot nack the message. Maybe connection closed. [endedWithCancellation]: %d", endedWithCancellation)
					panic(err)
				}

				log.Printf("Message nacked [messageId]: %s", d.MessageId)
			}

			log.Printf("consumer canceled the receiveing new messages, exit loop [queue]: %s - [consumer]: %s", source.QueueName, source.ConsumerTag)

			r.handlersDone.Done()
		}(source)
	}

	// wait for all handlers to finish processing all the messages
	r.handlersDone.Wait()
	log.Printf("Handlers ended")
	// wait for the whole consumer to finish all the stuff related to graceful shutdown before unblocking and exiting the method
	<-r.done
	return nil
}

// TearDown deals with graceful shutdown
// It's responsible to tell RabbitMQ to stop message sending
// Then it waits for handlers to finish messages already in process or in the buffer (prefetchSize)
// We cannot close connection before we ack/nack all the messages
// When handlers finished, we are safe to close connection
func (r *RabbitConsumer) TearDown() error {
	r.tearDownOnce.Do(func() {
		r.isTearDownInProgress = true
		defer func() {
			r.isTearDownInProgress = false
		}()

		for _, source := range r.sources {
			log.Printf("Cancel consumer: %s", source.ConsumerTag)

			if err := r.ch.Cancel(source.ConsumerTag, false); err != nil {
				log.Printf("consumer cancel failed: %s", err)
				continue
			}
		}

		log.Printf("Waiting for handlers to finish")
		r.handlersDone.Wait()

		log.Printf("close the connection")
		if err := r.conn.Close(); err != nil {
			log.Printf("attempt to close RabbitMQ connection failed: %s", err)
			return
		}

		defer log.Printf("Consumer teardown succesfully completed")

		r.isSuccessfully = true
		r.done <- true
	})

	return nil
}
