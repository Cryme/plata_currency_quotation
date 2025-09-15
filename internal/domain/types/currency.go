package types

import "encoding/json"

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	MXN Currency = "MXN"
)

func (c Currency) IsValid() bool {
	switch c {
	case USD, EUR, MXN:
		return true

	default:
		return false
	}
}

func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}
