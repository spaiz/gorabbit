job "orders_marketing" {
  datacenters = ["local"]

  group "orders_orders_marketing" {
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

    
        task "orders_marketing_0000" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0025 orders_marketing_0050"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0001" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0001 orders_marketing_0026"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0002" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0002 orders_marketing_0027"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0003" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0003 orders_marketing_0028"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0004" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0004 orders_marketing_0029"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0005" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0005 orders_marketing_0030"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0006" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0006 orders_marketing_0031"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0007" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0007 orders_marketing_0032"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0008" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0008 orders_marketing_0033"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0009" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0009 orders_marketing_0034"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0010" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0010 orders_marketing_0035"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0011" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0011 orders_marketing_0036"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0012" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0012 orders_marketing_0037"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0013" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0013 orders_marketing_0038"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0014" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0014 orders_marketing_0039"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0015" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0015 orders_marketing_0040"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0016" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0016 orders_marketing_0041"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0017" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0017 orders_marketing_0042"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0018" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0018 orders_marketing_0043"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0019" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0019 orders_marketing_0044"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0020" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0020 orders_marketing_0045"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0021" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0021 orders_marketing_0046"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0022" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0022 orders_marketing_0047"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0023" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0023 orders_marketing_0048"
            RABBITMQ_QUEUE_BINDING_KEY = "1"
            RABBITMQ_CONNECTION_STRING="amqp://user:bitnami@localhost:5672"
          }

          driver = "raw_exec"

          config {
            command = "consumer"
            args    = [""]
          }
        }
    
        task "orders_marketing_0024" {
        resources {
                cpu    = 20
                memory = 10
         }

          kill_timeout = "60s"

          env {
            IDENTITY = "orders_marketing"
            TOPIC = "orders"
            QUEUES = "orders_marketing_0024 orders_marketing_0049"
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