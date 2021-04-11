job "orders_processing" {
  datacenters = ["local"]

  group "orders_orders_processing" {
    restart {
      interval = "30m"
      attempts = 2
      delay    = "15s"
      mode     = "fail"
    }

    reschedule {
      attempts       = 1
      interval       = "24h"
      unlimited      = false
      delay          = "5s"
      delay_function = "constant"
    }

    
        task "orders_processing_0000" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0025 orders_processing_0050"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0001" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0001 orders_processing_0026"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0002" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0002 orders_processing_0027"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0003" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0003 orders_processing_0028"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0004" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0004 orders_processing_0029"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0005" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0005 orders_processing_0030"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0006" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0006 orders_processing_0031"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0007" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0007 orders_processing_0032"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0008" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0008 orders_processing_0033"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0009" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0009 orders_processing_0034"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0010" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0010 orders_processing_0035"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0011" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0011 orders_processing_0036"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0012" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0012 orders_processing_0037"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0013" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0013 orders_processing_0038"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0014" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0014 orders_processing_0039"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0015" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0015 orders_processing_0040"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0016" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0016 orders_processing_0041"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0017" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0017 orders_processing_0042"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0018" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0018 orders_processing_0043"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0019" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0019 orders_processing_0044"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0020" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0020 orders_processing_0045"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0021" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0021 orders_processing_0046"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0022" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0022 orders_processing_0047"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0023" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0023 orders_processing_0048"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_processing_0024" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_processing"
            TOPIC = "orders"
            QUEUES = "orders_processing_0024 orders_processing_0049"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
   }
}