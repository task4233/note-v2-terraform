package client

type OrderLog struct {
	Log Log `json:"log"`
}

type Log struct {
	Body string `json:"body"`
}
