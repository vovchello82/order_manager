package domain

import (
	"math/rand"
	"time"
)

type Order struct {
	Id              string
	PickUpLocation  string
	DropOffLocation string
	Customer        string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewRandomOrder() *Order {
	return &Order{
		Id:              randSeq(4),
		PickUpLocation:  randSeq(10),
		DropOffLocation: randSeq(10),
		Customer:        "testUser",
	}
}
