package Entities

type Instrument struct {
	ID             int
	Name           string
	Description    string
	Type           string
	Manufacturer   string
	Material       string
	Price          float64
	ImageURL       string
	StockQuantity  int
	Specifications map[string]any
}
