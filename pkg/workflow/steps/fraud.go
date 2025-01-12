package steps

import (
	"restate-order-demo/pkg/order"
	orderstatus "restate-order-demo/pkg/order/status"
	"restate-order-demo/pkg/services/fraud"

	restate "github.com/restatedev/sdk-go"
)

type FraudDetectedError struct{}

func (m *FraudDetectedError) Error() string {
	return "Fraud Detected"
}

func StartFraudCheck(ctx restate.Context, eb *EventBroker, custOrder *order.Order) error {
	eventError := EmitStatusUpdateEvent(ctx, eb, custOrder, orderstatus.PendingFraudReview, "Begin Fraud Review")

	if eventError != nil {
		return eventError
	}
	// Execute Fraud Check
	fraudOutput, fraudErr := restate.Service[fraud.FraudDecision](ctx, "Fraid", "ValidateOrder").Request(custOrder)

	if fraudErr != nil {
		return fraudErr
	}
	custOrder.RecoardFraudReviewDecision(fraudOutput.FraudDetected, fraudOutput.RejectionReason, fraudOutput.CheckDate)
	if fraudOutput.FraudDetected {
		//Emit fraud detected event
		eventError := EmitOrderStatusEvent(ctx, custOrder, orderstatus.Fraudlent, "Fraud Detected")
		if eventError != nil {
			return eventError
		}

		return &FraudDetectedError{}
		//TODO: What do we want to do at this point, cancel or have some intervention
	} else {
		//Emit Fraud Review Complete
		eventError := EmitOrderStatusEvent(ctx, custOrder, orderstatus.NoFraudDetected, "Order deemed not fraudlent")

		if eventError != nil {
			return eventError
		}
	}
	return nil
}
