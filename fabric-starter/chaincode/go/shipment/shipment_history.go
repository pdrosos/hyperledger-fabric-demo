package main

type ShipmentHistoryItem struct {
	TxId      string
	Value     string
	Timestamp string
	IsDelete  bool
}

type ShipmentHistory []ShipmentHistoryItem