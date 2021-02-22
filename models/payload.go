package models

type BasePayload struct {
	Status string
	Error  string
}

// Payload is the default response
type Payload struct {
	Payload interface{}
	BasePayload
}

// NewPayload populates and returns an initialised Payload struct.
func NewPayload(status string, payload interface{}, err string) Payload {
	return Payload{
		BasePayload: BasePayload{
			Status: status,
			Error:  err,
		},
		Payload: payload,
	}
}
