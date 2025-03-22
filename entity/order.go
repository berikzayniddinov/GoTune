package entity

type Order struct {
	ID              int64
	UserID          int64
	OrderDate       string
	DeliveryAddress string
	TotalPrice      float64
	Status          string
}
