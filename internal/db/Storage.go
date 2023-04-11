package domain

import "errors"

type Storage struct {
	db map[string]Order
}

func NewStorage() *Storage {
	return &Storage{
		db: make(map[string]Order, 10),
	}
}

func (s *Storage) PutOrder(o Order) error {
	s.db[o.Id] = o
	return nil
}

func (s *Storage) GetOrderById(id string) (*Order, error) {
	if o, exists := s.db[id]; exists {
		return &o, nil
	}
	return nil, errors.New("object no found")
}

func (s *Storage) DeleteOrderById(id string) {
	delete(s.db, id)
}

func (s *Storage) Size() int {
	return len(s.db)
}
