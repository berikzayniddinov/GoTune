package entity

type Order struct {
	ID              int
	UserID          int
	OrderDate       string
	DeliveryAddress string
	TotalPrice      float64
	Status          string
}
