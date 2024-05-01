package queue

import "encoding/json"

type TransferDto struct {
	ID      int     `json:"id"`
	PayerId int     `json:"payer_id"`
	PayeeId int     `json:"payee_id"`
	Value   float64 `json:"value"`
}

func (q *TransferDto) Marshal() ([]byte, error) {
	return json.Marshal(q)
}

func (q *TransferDto) Unmarhal(data []byte) error {
	return json.Unmarshal(data, q)
}
