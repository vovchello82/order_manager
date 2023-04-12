package domain

import (
	"errors"
	"log"
)

type Storage struct {
	db map[string]Order
}

func NewStorage() *Storage {
	return &Storage{
		db: make(map[string]Order, 10),
	}
}

func (s *Storage) PutOrder(o Order) {
	s.db[o.Id] = o
}

func (s *Storage) GetOrderById(id string) (*Order, error) {
	if o, exists := s.db[id]; exists {
		return &o, nil
	}
	return nil, errors.New("object no found")
}

func (s *Storage) DeleteOrderById(id string) {
	log.Printf("size before deleting %s : %d", id, len(s.db))
	delete(s.db, id)
	log.Printf("size after deleting %d", len(s.db))
}

func (s *Storage) Size() int {
	return len(s.db)
}
