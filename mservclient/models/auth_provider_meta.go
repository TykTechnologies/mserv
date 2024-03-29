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

// AuthProviderMeta auth provider meta
//
// swagger:model AuthProviderMeta
type AuthProviderMeta struct {

	// meta
	Meta interface{} `json:"meta,omitempty"`

	// name
	Name AuthProviderCode `json:"name,omitempty"`

	// storage engine
	StorageEngine StorageEngineCode `json:"storage_engine,omitempty"`
}

// Validate validates this auth provider meta
func (m *AuthProviderMeta) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStorageEngine(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AuthProviderMeta) validateName(formats strfmt.Registry) error {
	if swag.IsZero(m.Name) { // not required
		return nil
	}

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("name")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("name")
		}
		return err
	}

	return nil
}

func (m *AuthProviderMeta) validateStorageEngine(formats strfmt.Registry) error {
	if swag.IsZero(m.StorageEngine) { // not required
		return nil
	}

	if err := m.StorageEngine.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("storage_engine")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("storage_engine")
		}
		return err
	}

	return nil
}

// ContextValidate validate this auth provider meta based on the context it is used
func (m *AuthProviderMeta) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStorageEngine(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AuthProviderMeta) contextValidateName(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Name) { // not required
		return nil
	}

	if err := m.Name.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("name")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("name")
		}
		return err
	}

	return nil
}

func (m *AuthProviderMeta) contextValidateStorageEngine(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.StorageEngine) { // not required
		return nil
	}

	if err := m.StorageEngine.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("storage_engine")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("storage_engine")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AuthProviderMeta) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AuthProviderMeta) UnmarshalBinary(b []byte) error {
	var res AuthProviderMeta
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
