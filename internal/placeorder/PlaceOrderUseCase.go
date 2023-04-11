package placeorder

import (
	"encoding/json"
	domain "order_manager/internal/db"
	pub "order_manager/internal/publisher"
)

type PlaceOrderUseCase struct {
	publisher pub.Publisher
	db        *domain.Storage
}

func NewPlaceOrderUseCase(publisher pub.Publisher, storage *domain.Storage) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		publisher: publisher,
		db:        storage,
	}
}
func (puc *PlaceOrderUseCase) PlaceOrder(order domain.Order) error {
	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}
	puc.db.PutOrder(order)

	return puc.publisher.EmitObject(string(payload))
}
