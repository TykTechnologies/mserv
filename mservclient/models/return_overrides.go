// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ReturnOverrides return overrides
//
// swagger:model ReturnOverrides
type ReturnOverrides struct {

	// headers
	Headers map[string]string `json:"headers,omitempty"`

	// response code
	ResponseCode int32 `json:"response_code,omitempty"`

	// response error
	ResponseError string `json:"response_error,omitempty"`
}

// Validate validates this return overrides
func (m *ReturnOverrides) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this return overrides based on context it is used
func (m *ReturnOverrides) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ReturnOverrides) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReturnOverrides) UnmarshalBinary(b []byte) error {
	var res ReturnOverrides
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
