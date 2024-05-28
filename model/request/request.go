package request

import "encoding/json"

type Order struct {
	CustomerName  string           `json:"cus_name"`
	CustomerEmail string           `json:"cus_email"`
	Items         *json.RawMessage `json:"items"`
	Status        string           `json:"status"`
}
