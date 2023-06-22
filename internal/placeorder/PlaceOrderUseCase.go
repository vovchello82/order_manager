package placeorder

import (
	"context"
	"encoding/json"
	domain "order_manager/internal/db"
	pub "order_manager/internal/publisher"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
func (puc *PlaceOrderUseCase) PlaceOrder(ctx context.Context, order domain.Order) error {
	ctxCarrier := propagation.MapCarrier{}
	prop := otel.GetTextMapPropagator()
	prop.Inject(ctx, ctxCarrier)
	puc.db.PutOrder(order)

	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}
	out := map[string]interface{}{}
	json.Unmarshal([]byte(payload), &out)
	if err != nil {
		return err
	}

	for _, x := range ctxCarrier.Keys() {
		out[x] = ctxCarrier.Get(x)
	}

	payload, err = json.Marshal(out)
	if err != nil {
		return err
	}
	return puc.publisher.EmitObject(string(payload))
}
