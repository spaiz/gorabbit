package main

import (
	"fmt"
	"gorabbit/cmd/producer/app"
	"log"
	"os"
	"strconv"
)

func main() {
	connectionString := os.Getenv("RABBITMQ_CONNECTION_STRING")
	topic := os.Getenv("TOPIC")
	wd := os.Getenv("WORKING_DIRECTORY")

	delay, err := strconv.ParseInt(os.Getenv("DELAY"), 10, 64)
	if err != nil {
		log.Printf("DELAY was not provided, using default")
		delay = 1000
	}

	if wd != "" {
		log.Printf("Working directory from env: %s", wd)

		err := os.Chdir(wd)
		if err != nil {
			log.Fatalf("Failed to switch working directory: %s", err)
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}

	log.Printf("Current working direcotry: %s\n", dir)

	app := app.NewApp(
		app.RabbitConnectionString(connectionString),
		app.Topic(topic),
		app.ConnectionName(fmt.Sprintf("producer_%d", os.Getpid())),
		app.UseHeader(true),
		app.ExchangeType("fanout"),
		app.DataSetPath("./dataset/olist_orders_dataset.gz"),
		app.DataSetSkipHeader(true),
		app.Delay(delay),
		app.MaxReconnects(30),
		app.ReconnectDelay(2000),
		app.MaxSendMessageRetries(5),
	)

	log.Printf("Starting producer")
	err = app.Run()
	if err != nil {
		log.Fatalf("App finished with error: %s", err)
	}
}