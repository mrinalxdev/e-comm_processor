package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/workflow"
)

func OrderWorkflow(ctx workflow.Context, order Order) (string, error){
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	
	ctx = workflow.WithActivityOptions(ctx, options)
	
	var invResult, payResult, refundResult, restockResult ServiceResponse
	
	err := workflow.ExecuteActivity(ctx, CallNatsService, NatsSubjectPay, order).Get(ctx, &invResult)
	if err != nil || !invResult.Success {
		return "Payment failed", err
	}
	
	err = workflow.ExecuteActivity(ctx, CallNatsService, NatsSubjectInv, order).Get(ctx, &payResult)
	if err != nil || !payResult.Success {
		return "payment failed", err
	}
	
	
	//now we are just pretending that the order is waiting for shipping
	//so the new feature is defined as, checking if we should cancel
	
	if order.Item == "Golang Book"{
		fmt.Println("cancellation request recieved for 'Golang Book !!'")
		err = workflow.ExecuteActivity(ctx, CallNatsService, NatsSubjectRefund, order).Get(ctx, &refundResult)
        if err != nil {
            return "Refund Failed (Critical Error)", err
        }

        
        err = workflow.ExecuteActivity(ctx, CallNatsService, NatsSubjectRestock, order).Get(ctx, &restockResult)
        if err != nil {
            return "Restock Failed (Critical Error)", err
        }

        return "Order Cancelled & Refunded (Saga Complete)", nil
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