package models

type API struct {
	Request     string      `json:"request"`
	Transaction string      `json:"transaction,omitempty"`
	Body        interface{} `json:"body"`
}

type Example struct {
	Payload string `json:"payload"`
}

type ExampleTable struct {
	Seq uint64 `json:"seq"`
	Example
}
