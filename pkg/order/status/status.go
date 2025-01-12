package orderstatus

type OrderStatusCode string

type OrderStatus struct {
	Code   OrderStatusCode
	Reason string
}

const (
	Submitted OrderStatusCode = "submitted"

	PendingFraudReview OrderStatusCode = "pending-fraud"
	NoFraudDetected    OrderStatusCode = "no-fraud"
	Fraudlent          OrderStatusCode = "fraud"

	PendingCreditReview  OrderStatusCode = "pending-credit-review"
	CreditReviewApproved OrderStatusCode = "credit-review-approved"
	CreditReviewDenied   OrderStatusCode = "credit-review-denied"

	ApprovalRequired OrderStatusCode = "approval-required"
	Approved         OrderStatusCode = "approved"
	Rejected         OrderStatusCode = "rejected"
	Canceled         OrderStatusCode = "canceled"

	ReadyForFullfilment  OrderStatusCode = "ready-for-fulfillment"
	FullfilmentConfirmed OrderStatusCode = "fulfillment-confirmed"
)

func TerminalOrderStatus(status OrderStatusCode) bool {
	switch status {
	case Fraudlent, CreditReviewDenied, Rejected, Canceled:
		return true
	default:
		return false
	}
}
