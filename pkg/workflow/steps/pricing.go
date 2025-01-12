package steps

import (
	"log"
	"restate-order-demo/pkg/order"
	pricing "restate-order-demo/pkg/services"

	restate "github.com/restatedev/sdk-go"
)

func EnsureOrderIsPriced(ctx restate.Context, order *order.Order) error {
	// Execute Order Pricing
	if containsUnPriceLines(order) {
		log.Printf("Order requires pricing, pricing order.")
		pricedLines, priceErr := restate.Service[pricing.OrderLinePricing](ctx, "Pricing", "PriceOrder").Request(order)

		if priceErr != nil {
			return priceErr
		}
		updateOrderPricing(order, pricedLines)
	}
	return nil
}

func containsUnPriceLines(order *order.Order) bool {
	for _, line := range order.Lines {
		if line.Price == 0 {
			return true
		}
	}
	return false
}

func updateOrderPricing(order *order.Order, pricedLines pricing.OrderLinePricing) {
	for _, line := range order.Lines {
		if value, ok := pricedLines.LinePricing[line.Product]; ok {
			line.Price = value
		}
	}

}
