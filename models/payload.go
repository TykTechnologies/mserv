package models

type BasePayload struct {
	Status string
	Error  string
}

// Payload is the default response
type Payload struct {
	BasePayload
	Payload interface{}
}

func NewPayload(status string, payload interface{}, error string) Payload {
	return Payload{
		BasePayload: BasePayload{
			Status: status,
			Error:  error,
		},
		Payload: payload,
	}
}
