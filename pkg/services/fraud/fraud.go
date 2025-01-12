package fraud

import (
	"errors"
	"math/rand/v2"
	"restate-order-demo/pkg/order"
	"time"
)

type FraudDecision struct {
	FraudDetected   bool
	RejectionReason string
	CheckDate       time.Time
}

type Fraud struct{}

func (Fraud) ValidateOrder(order order.Order) (FraudDecision, error) {
	if rand.Float64() < 0.2 {
		return FraudDecision{}, errors.New("Something went wrong unknown why ;)")
	}

	result := FraudDecision{
		FraudDetected:   false,
		RejectionReason: "",
		CheckDate:       time.Now(),
	}
	if len(order.Lines) > 5 {
		result.FraudDetected = true
		result.RejectionReason = "Large Order"
	}
	return result, nil
}
