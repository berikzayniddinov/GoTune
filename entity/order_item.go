package entity

type OrderItem struct {
	ID           int64
	OrderID      int64
	InstrumentID int64
	Quantity     int
	Price        float64
}
