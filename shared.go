package main

const (
	TaskQueue = "ORDER_QUEUE"
	NatsSubjectInv = "service.inventory"
	NatsSubjectPay = "service.payment"
	NatsSubjectRefund = "service.refund"
	NatsSubjectRestock = "service.restock"
)

type Order struct {
	OrderID string
	Amount float64
	Item string
}

type ServiceResponse struct {
	Success bool
	Message string
}