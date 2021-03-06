// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Object object
//
// swagger:model Object
type Object struct {

	// hook name
	HookName string `json:"hook_name,omitempty"`

	// metadata
	Metadata map[string]string `json:"metadata,omitempty"`

	// spec
	Spec map[string]string `json:"spec,omitempty"`

	// hook type
	HookType HookType `json:"hook_type,omitempty"`

	// request
	Request *MiniRequestObject `json:"request,omitempty"`

	// session
	Session *SessionState `json:"session,omitempty"`
}

// Validate validates this object
func (m *Object) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHookType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRequest(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSession(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Object) validateHookType(formats strfmt.Registry) error {

	if swag.IsZero(m.HookType) { // not required
		return nil
	}

	if err := m.HookType.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("hook_type")
		}
		return err
	}

	return nil
}

func (m *Object) validateRequest(formats strfmt.Registry) error {

	if swag.IsZero(m.Request) { // not required
		return nil
	}

	if m.Request != nil {
		if err := m.Request.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("request")
			}
			return err
		}
	}

	return nil
}

func (m *Object) validateSession(formats strfmt.Registry) error {

	if swag.IsZero(m.Session) { // not required
		return nil
	}

	if m.Session != nil {
		if err := m.Session.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("session")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Object) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Object) UnmarshalBinary(b []byte) error {
	var res Object
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
