package app

// Option defines abstraction for variadic functions used to config the behaviour
type Option func(s *App)

// RabbitConnectionString sets connection string to rabbitMQ
func RabbitConnectionString(connectionString string) func(app *App) {
	return func(a *App) {
		a.config.RabbitMQConnectionString = connectionString
	}
}

// Queues sets queue names that consumer will subscribe to
func Queues(queues []string) func(app *App) {
	return func(a *App) {
		a.config.QueueNames = queues
	}
}

// Identity will be displayed in RabbitMQ management UI
func Identity(identity string) func(app *App) {
	return func(a *App) {
		a.config.Identity = identity
	}
}

// UseHeader will be displayed in RabbitMQ management UI
func UseHeader(useHeader bool) func(app *App) {
	return func(a *App) {
		a.config.UseHeader = useHeader
	}
}

// Topic will be displayed in RabbitMQ management UI
func Topic(topic string) func(app *App) {
	return func(a *App) {
		a.config.Topic = topic
	}
}

// SingleActiveConsumer uses rabbitMQ's ability to ensure only single consumer will receive messages
func SingleActiveConsumer(enabled bool) func(app *App) {
	return func(a *App) {
		a.config.SingleActiveConsumer = enabled
	}
}

// BindingKey will used for binding queues to exchange
func BindingKey(key string) func(app *App) {
	return func(a *App) {
		a.config.BindingKey = key
	}
}

// prefetchSize defines the amount fo messages pushed by rabbitMQ and kept by the driver
// locally in the buffer. Be careful with it.
// It can be useful but also can create issues when size it too high.
// It can create busy consumer y consuming messages and keeping it in the buffer,
// when other consumers will be idle
func PrefetchSize(size int) func(app *App) {
	return func(a *App) {
		a.config.PrefetchSize = size
	}
}