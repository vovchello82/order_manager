package placeorder

import (
	"encoding/json"
	domain "order_manager/internal/db"
	pub "order_manager/internal/publisher"
)

type PlaceOrderUseCase struct {
	publisher pub.Publisher
}

func NewPlaceOrderUseCase(publisher pub.Publisher) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		publisher: publisher,
	}
}
func (puc *PlaceOrderUseCase) PlaceOrder(order domain.Order) error {
	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return puc.publisher.EmitObject(string(payload))
}
