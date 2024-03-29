// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ValidatePathMeta validate path meta
//
// swagger:model ValidatePathMeta
type ValidatePathMeta struct {

	// Allows override of default 422 Unprocessible Entity response code for validation errors.
	ErrorResponseCode int64 `json:"error_response_code,omitempty"`

	// method
	Method string `json:"method,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// schema
	Schema interface{} `json:"schema,omitempty"`

	// schema b64
	SchemaB64 string `json:"schema_b64,omitempty"`
}

// Validate validates this validate path meta
func (m *ValidatePathMeta) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this validate path meta based on context it is used
func (m *ValidatePathMeta) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ValidatePathMeta) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ValidatePathMeta) UnmarshalBinary(b []byte) error {
	var res ValidatePathMeta
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
