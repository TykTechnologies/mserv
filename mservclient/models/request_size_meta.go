// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RequestSizeMeta request size meta
//
// swagger:model RequestSizeMeta
type RequestSizeMeta struct {

	// method
	Method string `json:"method,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// size limit
	SizeLimit int64 `json:"size_limit,omitempty"`
}

// Validate validates this request size meta
func (m *RequestSizeMeta) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this request size meta based on context it is used
func (m *RequestSizeMeta) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RequestSizeMeta) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RequestSizeMeta) UnmarshalBinary(b []byte) error {
	var res RequestSizeMeta
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
