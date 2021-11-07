package client

type Order struct {
	// commented out because this field is not used for now
	// ID    string      `json:"id"`
	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	Log Log `json:"log"`
}

type Log struct {
	Body string `json:"body"`
}
