// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RoutingTrigger routing trigger
//
// swagger:model RoutingTrigger
type RoutingTrigger struct {

	// rewrite to
	RewriteTo string `json:"rewrite_to,omitempty"`

	// on
	On RoutingTriggerOnType `json:"on,omitempty"`

	// options
	Options *RoutingTriggerOptions `json:"options,omitempty"`
}

// Validate validates this routing trigger
func (m *RoutingTrigger) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateOn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOptions(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RoutingTrigger) validateOn(formats strfmt.Registry) error {
	if swag.IsZero(m.On) { // not required
		return nil
	}

	if err := m.On.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("on")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("on")
		}
		return err
	}

	return nil
}

func (m *RoutingTrigger) validateOptions(formats strfmt.Registry) error {
	if swag.IsZero(m.Options) { // not required
		return nil
	}

	if m.Options != nil {
		if err := m.Options.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("options")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("options")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this routing trigger based on the context it is used
func (m *RoutingTrigger) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateOn(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateOptions(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RoutingTrigger) contextValidateOn(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.On) { // not required
		return nil
	}

	if err := m.On.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("on")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("on")
		}
		return err
	}

	return nil
}

func (m *RoutingTrigger) contextValidateOptions(ctx context.Context, formats strfmt.Registry) error {

	if m.Options != nil {

		if swag.IsZero(m.Options) { // not required
			return nil
		}

		if err := m.Options.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("options")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("options")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RoutingTrigger) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RoutingTrigger) UnmarshalBinary(b []byte) error {
	var res RoutingTrigger
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
