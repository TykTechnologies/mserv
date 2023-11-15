// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HardTimeoutMeta hard timeout meta
//
// swagger:model HardTimeoutMeta
type HardTimeoutMeta struct {

	// method
	Method string `json:"method,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// time out
	TimeOut int64 `json:"timeout,omitempty"`
}

// Validate validates this hard timeout meta
func (m *HardTimeoutMeta) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hard timeout meta based on context it is used
func (m *HardTimeoutMeta) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HardTimeoutMeta) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HardTimeoutMeta) UnmarshalBinary(b []byte) error {
	var res HardTimeoutMeta
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}