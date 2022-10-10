package models

// default Response structure
//
type R struct {
	Status      uint16      `json:"status"`
	Transaction string      `json:"transaction,omitempty"`
	Response    interface{} `json:"response"`
}
