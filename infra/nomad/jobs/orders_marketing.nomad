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
            QUEUES = "orders_marketing_0002 orders_marketing_0004"
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
            QUEUES = "orders_marketing_0001 orders_marketing_0003"
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