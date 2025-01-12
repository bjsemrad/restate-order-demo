package order

import (
	orderstatus "restate-order-demo/pkg/order/status"
	"time"
)

type Order struct {
	Channel                string
	Status                 *orderstatus.OrderStatus
	OrderNumber            string
	FullfilmentOrderNumber string `json:"FullfilmentOrderNumber,omitempty"`
	Lines                  []*OrderLine
	Payment                *Payment
	LastUpdated            time.Time
	Settings               *OrderSettings
	PipelineMetadata       *OrderPipelineMetadata
}

type OrderSettings struct {
	LiftGateRequired     bool
	PackingListInEachBox bool
}

type OrderLine struct {
	LineNumber int
	Product    string
	Quantity   int
	Price      float64
}

type Payment struct {
	CreditCard    string `json:"CreditCard,omitempty"`
	AccountNumber string `json:"AccountNumber,omitempty"`
}

type OrderPipelineMetadata struct {
	FraudReview   *OrderFraudReview     `json:"FraudReview,omitempty"`
	CreditReview  *OrderCreditReview    `json:"CreditReview,omitempty"`
	StatusHistory []*OrderStatusHistory `json:"StatusHistory,omitempty"`
}

type OrderStatusHistory struct {
	Status orderstatus.OrderStatusCode
	Reason string
	Date   time.Time
}

type OrderFraudReview struct {
	FraudDetected   bool
	RejectionReason string
	DecisionDate    time.Time
}

type OrderCreditReview struct {
	CreditAvailable bool      `json:"CreditAvailable,omitempty"`
	AvailableCredit float64   `json:"AvailableCredit,omitempty"`
	NewLimit        float64   `json:"NewLimit,omitempty"`
	CreditDecision  string    `json:"CreditDecision,omitempty"`
	Reviewier       string    `json:"Reviewer,omitempty"`
	DecisionDate    time.Time `json:"DecisionDate,omitempty"`
}

func NewOrder(orderNumber string) *Order {
	return &Order{
		OrderNumber:      orderNumber,
		Lines:            []*OrderLine{},
		Payment:          &Payment{},
		PipelineMetadata: &OrderPipelineMetadata{},
	}
}

func (o *Order) AddLine(lineNumber int, product string, quantity int, price float64) {
	line := &OrderLine{
		LineNumber: lineNumber,
		Product:    product,
		Quantity:   quantity,
		Price:      price,
	}
	o.Lines = append(o.Lines, line)
}

func (o *Order) UpdateStatus(newStatus orderstatus.OrderStatusCode, reason string) {
	o.recordStatusChange(o.Status.Code, o.Status.Reason, o.LastUpdated)
	o.Status.Code = newStatus
	o.Status.Reason = reason
	o.LastUpdated = time.Now()
}

func (o *Order) Total() float64 {
	total := 0.0
	for _, line := range o.Lines {
		total += line.Price * float64(line.Quantity)
	}
	return total
}

func (o *Order) RecordCreditReservation(creditAvailable bool, availableCredit float64) {
	o.ensureMetadataInitalized()
	o.PipelineMetadata.CreditReview = &OrderCreditReview{
		CreditAvailable: creditAvailable,
		AvailableCredit: availableCredit,
	}
}

func (o *Order) RecordCreditReviewDecision(decision string, reviewer string, newLimit float64, reviewDate time.Time) {
	o.ensureMetadataInitalized()
	o.PipelineMetadata.CreditReview = &OrderCreditReview{
		CreditDecision: decision,
		Reviewier:      reviewer,
		DecisionDate:   reviewDate,
		NewLimit:       newLimit,
	}
}

func (o *Order) RecoardFraudReviewDecision(fradulent bool, reason string, reviewDate time.Time) {
	o.ensureMetadataInitalized()
	o.PipelineMetadata.FraudReview = &OrderFraudReview{
		FraudDetected:   fradulent,
		RejectionReason: reason,
		DecisionDate:    reviewDate,
	}
}

func (o *Order) recordStatusChange(status orderstatus.OrderStatusCode, reason string, statusDate time.Time) {
	o.ensureMetadataInitalized()
	o.PipelineMetadata.StatusHistory = append(o.PipelineMetadata.StatusHistory, &OrderStatusHistory{
		Status: status,
		Reason: reason,
		Date:   statusDate,
	})
}

func (o *Order) ensureMetadataInitalized() {
	if o.PipelineMetadata == nil {
		o.PipelineMetadata = &OrderPipelineMetadata{
			StatusHistory: make([]*OrderStatusHistory, 0),
		}
	}
}
