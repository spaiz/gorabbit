package main

import (
	"gorabbit/cmd/consumer/app"
	"log"
	"os"
	"strings"
)

func main() {
	identity := os.Getenv("IDENTITY")
	topic := os.Getenv("TOPIC")
	queuesArg := os.Getenv("QUEUES")
	connectionString := os.Getenv("RABBITMQ_CONNECTION_STRING")
	bindingKey := os.Getenv("RABBITMQ_QUEUE_BINDING_KEY")

	queues := extractQueues(queuesArg)

	app := app.NewApp(
		app.Topic(topic),
		app.RabbitConnectionString(connectionString),
		app.Queues(queues),
		app.Identity(identity),
		app.SingleActiveConsumer(true),
		app.BindingKey(bindingKey),
		app.PrefetchSize(1),
		app.UseHeader(true),
	)

	log.Printf("Starting consumer")
	err := app.Run()
	if err != nil {
		log.Fatalf("App finished with error: %s", err)
	}
}

func extractQueues(queues string) []string {
	return strings.Split(queues, " ")
}
