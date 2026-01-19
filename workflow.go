package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/workflow"
)

func OrderWorkflow(ctx workflow.Context, order Order) (string, error){
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	
	ctx = workflow.WithActivityOptions(ctx, options)
	
	var invResult ServiceResponse
	var payResult ServiceResponse
	
	err := workflow.ExecuteActivity(ctx, CallNatsService, NatsSubjectPay, order).Get(ctx, &invResult)
	if err != nil || !invResult.Success {
		return "Payment failed", err
	}
	
	err = workflow.ExecuteActivity(ctx, CallNatsService, NatsSubjectInv, order).Get(ctx, &payResult)
	if err != nil || !payResult.Success {
		return "payment failed", err
	}
	
	return "Order Processed Successfully", nil
}


func CallNatsService(subject string, order Order) (ServiceResponse, error) {
	nc, _ := nats.Connect(nats.DefaultURL)
	
	
	defer nc.Close()
	reqData, _ := json.Marshal(order)
	msg, err := nc.Request(subject, reqData, 2*time.Second)
	if err != nil {
		return ServiceResponse{}, err
	}
	
	var resp ServiceResponse
	json.Unmarshal(msg.Data, &resp)
	
	if !resp.Success {
		return resp, errors.New(resp.Message)
	}
	
	return resp, nil
}