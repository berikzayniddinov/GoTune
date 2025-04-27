package events

type CartCheckedOutEvent struct {
	UserID string           `json:"user_id"`
	Items  []CheckedOutItem `json:"items"`
}

type CheckedOutItem struct {
	InstrumentID string `json:"instrument_id"`
	Quantity     int32  `json:"quantity"`
}
