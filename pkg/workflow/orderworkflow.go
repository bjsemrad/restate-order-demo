package orderworkflow

import (
	"errors"
	"log"
	"restate-order-demo/pkg/order"
	"restate-order-demo/pkg/workflow/steps"

	restate "github.com/restatedev/sdk-go"
)

type OrderWorkflow struct {
	eventProducter *steps.EventBroker
}

func InitializeWorkflow(eventProducter *steps.EventBroker) *OrderWorkflow {
	return &OrderWorkflow{
		eventProducter: eventProducter,
	}
}

func (w *OrderWorkflow) ProcessOrder(ctx restate.Context, custOrder *order.Order) (order.Order, error) { //TODO comeback to the workflow outpu

	log.Printf("Start Ensuring Order is Priced: %s.\n", custOrder.OrderNumber)
	steps.EnsureOrderIsPriced(ctx, custOrder)

	log.Printf("Start Fraud Check: %s.\n", custOrder.OrderNumber)
	fraudErr := steps.StartFraudCheck(ctx, w.eventProducter, custOrder)
	var fraudDetectedError *steps.FraudDetectedError
	if fraudErr != nil && !errors.As(fraudErr, &fraudDetectedError) {
		return *custOrder, fraudErr
	}
	log.Printf("End Fraud Check Step: %s.\n", custOrder.OrderNumber)
	//
	// log.Printf("Start Validate Business Rules Check Step: %s.\n", custOrder.OrderNumber)
	// //TODO: Validate Business Rules
	// log.Printf("End Validate Business Rules Step: %s.\n", custOrder.OrderNumber)
	//
	// log.Printf("Start Cust Approval Step: %s.\n", custOrder.OrderNumber)
	// //TODO: Cust Approval
	// log.Printf("End Cust Approval Step: %s.\n", custOrder.OrderNumber)
	//
	// if custOrder.Payment != nil && strings.Trim(custOrder.Payment.AccountNumber, " ") != "" {
	// 	log.Printf("Start Credit Review Step: %s.\n", custOrder.OrderNumber)
	// 	creditReviewErr := orderworkflowstep.StartCreditReview(ctx, custOrder)
	// 	var creditDeniedError *orderworkflowstep.CreditDeniedError
	// 	if creditReviewErr != nil && !errors.As(creditReviewErr, &creditDeniedError) {
	// 		return *custOrder, creditReviewErr
	// 	}
	// 	log.Printf("End Credit Review Step: %s.\n", custOrder.OrderNumber)
	//
	// }
	//
	// log.Printf("Start Apply Setting Step: %s.\n", custOrder.OrderNumber)
	// settingsErr := orderworkflowstep.ApplySettings(ctx, custOrder)
	// if settingsErr != nil {
	// 	return *custOrder, settingsErr
	// }
	// log.Printf("End Apply Setting Step: %s.\n", custOrder.OrderNumber)
	//
	// log.Printf("Start Ops Rules Step: %s.\n", custOrder.OrderNumber)
	// //TODO: Operational Rule Checks
	// log.Printf("End Ops Rules Step: %s.\n", custOrder.OrderNumber)
	//
	// log.Printf("Start Intervention Step: %s.\n", custOrder.OrderNumber)
	// //TODO: Team Intervention
	// log.Printf("End Intervention Step: %s.\n", custOrder.OrderNumber)
	//
	// //IF we are not in a terminal status send the fulfillment signal
	// if !orderstatus.TerminalOrderStatus(custOrder.Status.Code) {
	// 	log.Printf("Start Prepare For Fulfillment Step: %s.\n", custOrder.OrderNumber)
	// 	fulfillError := orderworkflowstep.PrepareOrderForFulfillment(ctx, custOrder)
	// 	if fulfillError != nil {
	// 		return *custOrder, fulfillError
	// 	}
	// 	log.Printf("End Prepare for Fulfillment Step: %s.\n", custOrder.OrderNumber)
	//
	// 	log.Printf("Start Wait for Confirmed Order: %s.\n", custOrder.OrderNumber)
	// 	orderworkflowstep.WaitForConfirmedOrder(ctx, custOrder)
	// }
	log.Printf("Workflow Complete: %s.\n", custOrder.OrderNumber)
	return *custOrder, nil
}
