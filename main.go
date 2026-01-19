package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	mode := flag.String("mode", "worker", "Mode: 'worker' or 'trigger'")
	flag.Parse()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	if *mode == "worker" {
		nc, _ := nats.Connect(nats.DefaultURL)
		StartNatsServices(nc)

		w := worker.New(c, TaskQueue, worker.Options{})
		w.RegisterWorkflow(OrderWorkflow)
		w.RegisterActivity(CallNatsService)

		log.Println("worker started. Waiting for tasks")
		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("Unable to start worker", err)
		}

	} else {
		order := Order{OrderID: "101", Amount: 99.50, Item: "Golang Book"}
		we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
			ID:        "order-101",
			TaskQueue: TaskQueue,
		}, OrderWorkflow, order)
		
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		
		log.Printf("workflow started (ID: %sRunID: %s)", we.GetID(), we.GetRunID())
	
		var result string
		we.Get(context.Background(), &result)
		fmt.Printf("workflow Result: %s\n", result)
	}
}