// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
)

// AccessRequestType AccessRequestType is the type for OAuth param `grant_type`
//
// swagger:model AccessRequestType
type AccessRequestType string

// Validate validates this access request type
func (m AccessRequestType) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this access request type based on context it is used
func (m AccessRequestType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}