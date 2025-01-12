package pricing

import (
	"errors"
	"math/rand/v2"
	"restate-order-demo/pkg/order"

	restate "github.com/restatedev/sdk-go"
)

type OrderLinePricing struct {
	LinePricing map[string]float64
}

type Pricing struct{}

func (Pricing) PriceOrder(ctx restate.Context, order *order.Order) (OrderLinePricing, error) {
	if rand.Float64() < 0.4 {
		return OrderLinePricing{}, errors.New("I felt like it")
	}
	result := OrderLinePricing{
		LinePricing: make(map[string]float64, 0),
	}
	for _, line := range order.Lines {
		if line.Price == 0 {
			result.LinePricing[line.Product] = 1 + rand.Float64()*(500-1)
		}
	}
	return result, nil
}
