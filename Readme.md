# POC for playing with RabbitMQ and Constant Hash Exchange Plugin

**Pre requirements**

- MacOS
- Docker installed
- go1.14.1
- Nomad installed

## Quick Start

Start Nomad agent

```bash
nomad agent -dev -bind 0.0.0.0 -log-level DEBUG -dc local
```
Check Nomad is up by visiting: http://localhost:4646/ui/jobs

All the commands hsould be run from this directiry:

```bash
path_to_repository/gorabbit/infra/cli

```

Run this command to start RabbitMQ

```bash
go run cli.go SystemUp

```
It will start the Nomad job and download RabbitMQ image and run it.

Visit the RabbitMQ UI: http://127.0.0.1:15672/

> Username: user
Password: bitnami

Then install the consumers:

```bash
go run cli.go install

```

Run consumers

```bash
go run cli.go up

```

New jobs in Nomad will be started. You will see it in UI.

Run producer (change path to the repository)

```bash
export RABBITMQ_CONNECTION_STRING=amqp://user:bitnami@localhost:5672
export WORKING_DIRECTORY=/path_to_repository/gorabbit
export TOPIC=orders producer
producer
```

## CLI

Run these commands form the ../gorabbit/infra/cli directory, cause internal paths are relative to it

Build and install producer and consumer globally

```bash
go run cli.go install
```

To generate new mappings from the *infra/cli/manifest.json* file run:
```bash
go run cli.go remap

```
Start RabbitMQ

```bash
go run cli.go systemUp

```
Start consumers

```bash
go run cli.go up

```
Then, you can change the instance number inside the *infra/cli/manifest.json* , and rescale.

```bash
go run cli.go rescale

```
This command will regenerate the mappings, update the Nomad job files and restart the jobs.

For other available commands look to the file *infra/cli/magefile.go*