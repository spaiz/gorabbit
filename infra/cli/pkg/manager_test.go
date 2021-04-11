package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseManifest(t *testing.T) {
	assert := assert.New(t)

	filePath := "manifest.json"
	manifest, err := LoadManifest(filePath)
	assert.Nilf(err, "Message: %s", err)
	expected := &Manifest{
		Workers: []*Worker{
			{
				Identity:                 "orders_processing",
				Topic:                    "orders",
				QueuesSet:                QueuesSet{
					PartitionsNumber: 10,
				},
				InstancesNumber:          10,
				RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
			},
			{
				Identity:                 "orders_marketing",
				Topic:                    "orders",
				QueuesSet:                QueuesSet{
					PartitionsNumber: 10,
				},
				InstancesNumber:          5,
				RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
			},
		},
	}
	assert.Equal(expected, manifest, "they should be equal")
}

func Test_CreateState(t *testing.T) {
	assert := assert.New(t)

	manifest := &Manifest{
		Workers: []*Worker{
			{
				Identity:                 "orders_processing",
				Topic:                    "orders",
				QueuesSet:                QueuesSet{
					PartitionsNumber: 10,
				},
				InstancesNumber:          10,
				RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
			},
			{
				Identity:                 "orders_marketing",
				Topic:                    "orders",
				QueuesSet:                QueuesSet{
					PartitionsNumber: 10,
				},
				InstancesNumber:          5,
				RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
			},
		},
	}

	state, err := manifest.CreateState()
	assert.Nilf(err, "Message: %s", err)

	expected := &State{
		mappings: []*QueuesMapping{
			{
				instances: []*WorkerInstance{
					{
						Identity:                 "orders_processing",
						InstanceId:               1,
						queues:                   []string{"orders_processing_0010"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               2,
						queues:                   []string{"orders_processing_0001"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               3,
						queues:                   []string{"orders_processing_0002"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               4,
						queues:                   []string{"orders_processing_0003"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               5,
						queues:                   []string{"orders_processing_0004"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               6,
						queues:                   []string{"orders_processing_0005"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               7,
						queues:                   []string{"orders_processing_0006"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               8,
						queues:                   []string{"orders_processing_0007"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               9,
						queues:                   []string{"orders_processing_0008"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_processing",
						InstanceId:               10,
						queues:                   []string{"orders_processing_0009"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
				},
				topic:     "orders",
				identity:  "orders_processing",
			},
			{
				instances: []*WorkerInstance{
					{
						Identity:                 "orders_marketing",
						InstanceId:               1,
						queues:                   []string{"orders_marketing_0005", "orders_marketing_0010"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_marketing",
						InstanceId:               2,
						queues:                   []string{"orders_marketing_0001", "orders_marketing_0006"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_marketing",
						InstanceId:               3,
						queues:                   []string{"orders_marketing_0002", "orders_marketing_0007"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_marketing",
						InstanceId:               4,
						queues:                   []string{"orders_marketing_0003", "orders_marketing_0008"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
					{
						Identity:                 "orders_marketing",
						InstanceId:               5,
						queues:                   []string{"orders_marketing_0004", "orders_marketing_0009"},
						RabbitMQConnectionString: "amqp://user:bitnami@localhost:5672",
					},
				},
				topic:     "orders",
				identity:  "orders_marketing",
			},
		},
	}

	assert.Equal(expected, state, "they should be equal")
}

func Test_CreateJob(t *testing.T) {
	t.Skip("Un-skip if you want play with that")
	assert := assert.New(t)

	manifest, err := LoadManifest("manifest.json")
	assert.Nil(err)

	state, err := manifest.CreateState()
	assert.Nil(err)

	jobs, err := manifest.CreateJobs(state)
	assert.Nil(err)

	err = manifest.DumpJobs(jobs, "../nomad/jobs/")
	assert.Nil(err)
}