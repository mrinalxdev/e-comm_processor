package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func StartNatsServices(nc *nats.Conn) {
	nc.Subscribe(NatsSubjectInv, func(m *nats.Msg) {
		var order Order
		json.Unmarshal(m.Data, &order)
		
		fmt.Printf("[NATS Service] Checking Inventory for: %s\n", order.Item)
		
		resp := ServiceResponse{Success: true, Message: "Item Reserved"}
		data, _ := json.Marshal(resp)
		m.Respond(data) 
	})

	nc.Subscribe(NatsSubjectPay, func(m *nats.Msg) {
		var order Order
		json.Unmarshal(m.Data, &order)

		fmt.Printf("[NATS Service] Charging Credit Card: $%.2f\n", order.Amount)

		resp := ServiceResponse{Success: true, Message: "Transaction Approved"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})

	log.Println("nats Microservices (Inventory & Payment) are running...")
}