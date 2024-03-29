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

// EventHandlerTriggerConfig event handler trigger config
//
// swagger:model EventHandlerTriggerConfig
type EventHandlerTriggerConfig struct {

	// handler meta
	HandlerMeta interface{} `json:"handler_meta,omitempty"`

	// handler name
	HandlerName TykEventHandlerName `json:"handler_name,omitempty"`
}

// Validate validates this event handler trigger config
func (m *EventHandlerTriggerConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHandlerName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EventHandlerTriggerConfig) validateHandlerName(formats strfmt.Registry) error {
	if swag.IsZero(m.HandlerName) { // not required
		return nil
	}

	if err := m.HandlerName.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("handler_name")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("handler_name")
		}
		return err
	}

	return nil
}

// ContextValidate validate this event handler trigger config based on the context it is used
func (m *EventHandlerTriggerConfig) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateHandlerName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EventHandlerTriggerConfig) contextValidateHandlerName(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.HandlerName) { // not required
		return nil
	}

	if err := m.HandlerName.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("handler_name")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("handler_name")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *EventHandlerTriggerConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EventHandlerTriggerConfig) UnmarshalBinary(b []byte) error {
	var res EventHandlerTriggerConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
