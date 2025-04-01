package Entities

type OrderItem struct {
	ID           int
	OrderID      int
	InstrumentID int
	Quantity     int
	Price        float64
}
