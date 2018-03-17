package inputmodel

type ShipmentHistoryRecord struct {
	TxId      string `json:"txId"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
}

type ShipmentHistory []ShipmentHistoryRecord
