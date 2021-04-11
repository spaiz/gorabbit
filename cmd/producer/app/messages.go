package app

import "time"

// OrderMessage represents our business message
type OrderMessage struct {
	OrderId string
	CustomerId string
	MessageType string
	PurchaseTimestamp *time.Time
	ApprovedAt *time.Time
	DeliveredCarrierDate *time.Time
	DeliveredCustomerDate *time.Time
	EstimatedDeliveryDate *time.Time
}