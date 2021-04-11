package app

import "time"

// Config to configure behaviour
type Config struct {
	RabbitMQConnectionString string
	exchangeType             string
	connectionName           string
	topic                    string
	useHeader                bool
	dataSetPath              string
	dataSetSkipHeader        bool  // skip first line in csv file
	Delay                    int64 // Milliseconds
	MaxReconnects            int   // reconnect attempts on send message
	ReconnectDelay           time.Duration
	MaxSendMessageRetries    int
}

// Option defines abstraction for variadic functions used to config the behaviour
type Option func(s *App)

// RabbitConnectionString sets connection string to rabbitMQ
func RabbitConnectionString(connectionString string) func(app *App) {
	return func(a *App) {
		a.config.RabbitMQConnectionString = connectionString
	}
}

// topic sets exchange name to be created when producer will start up
func Topic(topic string) func(app *App) {
	return func(a *App) {
		a.config.topic = topic
	}
}

// ConnectionName sets connection name to be displayed at RabbitMQ management UI
func ConnectionName(connectionName string) func(app *App) {
	return func(a *App) {
		a.config.connectionName = connectionName
	}
}

// UseHeader forces producer to use header for routing key
func UseHeader(useHeader bool) func(app *App) {
	return func(a *App) {
		a.config.useHeader = useHeader
	}
}

// ExchangeType sets the type of the exchange
// Once created, cannot be changed without deleting it firstly from rabbitMQ
func ExchangeType(exchangeType string) func(app *App) {
	return func(a *App) {
		a.config.exchangeType = exchangeType
	}
}

// DataSetPath sets path to file with csv that will be used as test data to be sent to RabbitMQ
func DataSetPath(path string) func(app *App) {
	return func(a *App) {
		a.config.dataSetPath = path
	}
}

// DataSetSkipHeader skip first header line in csv file
func DataSetSkipHeader(skip bool) func(app *App) {
	return func(a *App) {
		a.config.dataSetSkipHeader = skip
	}
}

// Delay defines delay between each message sending in milliseconds
func Delay(d int64) func(app *App) {
	return func(a *App) {
		a.config.Delay = d
	}
}

// MaxReconnects defines reconnect attempts on when send message fails
func MaxReconnects(attempts int) func(app *App) {
	return func(a *App) {
		a.config.MaxReconnects = attempts
	}
}

// ReconnectDelay defines delay between reconnects attempts in millisecinds
func ReconnectDelay(d int64) func(app *App) {
	return func(a *App) {
		a.config.ReconnectDelay = time.Millisecond * time.Duration(d)
	}
}

// MaxSendMessageRetries defines max send message retries on connection closed error
func MaxSendMessageRetries(retries int) func(app *App) {
	return func(a *App) {
		a.config.MaxSendMessageRetries = retries
	}
}