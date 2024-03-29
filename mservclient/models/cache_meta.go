// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CacheMeta cache meta
//
// swagger:model CacheMeta
type CacheMeta struct {

	// cache key regex
	CacheKeyRegex string `json:"cache_key_regex,omitempty"`

	// method
	Method string `json:"method,omitempty"`

	// path
	Path string `json:"path,omitempty"`
}

// Validate validates this cache meta
func (m *CacheMeta) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this cache meta based on context it is used
func (m *CacheMeta) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CacheMeta) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CacheMeta) UnmarshalBinary(b []byte) error {
	var res CacheMeta
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
