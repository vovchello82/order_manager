package placeorder

import (
	domain "order_manager/internal/db"
)

type PlaceOrder interface {
	CreateOrder(order domain.Order) error
}
