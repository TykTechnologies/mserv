// Code generated by go-swagger; DO NOT EDIT.

package mw

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/TykTechnologies/mserv/mservclient/models"
)

// MwUpdateReader is a Reader for the MwUpdate structure.
type MwUpdateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MwUpdateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewMwUpdateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewMwUpdateInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[PUT /api/mw/{id}] mwUpdate", response, response.Code())
	}
}

// NewMwUpdateOK creates a MwUpdateOK with default headers values
func NewMwUpdateOK() *MwUpdateOK {
	return &MwUpdateOK{}
}

/*
MwUpdateOK describes a response with status code 200, with default header values.

Response that only includes the ID of the middleware as `BundleID` in the `Payload`
*/
type MwUpdateOK struct {
	Payload *MwUpdateOKBody
}

// IsSuccess returns true when this mw update o k response has a 2xx status code
func (o *MwUpdateOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this mw update o k response has a 3xx status code
func (o *MwUpdateOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this mw update o k response has a 4xx status code
func (o *MwUpdateOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this mw update o k response has a 5xx status code
func (o *MwUpdateOK) IsServerError() bool {
	return false
}

// IsCode returns true when this mw update o k response a status code equal to that given
func (o *MwUpdateOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the mw update o k response
func (o *MwUpdateOK) Code() int {
	return 200
}

func (o *MwUpdateOK) Error() string {
	return fmt.Sprintf("[PUT /api/mw/{id}][%d] mwUpdateOK  %+v", 200, o.Payload)
}

func (o *MwUpdateOK) String() string {
	return fmt.Sprintf("[PUT /api/mw/{id}][%d] mwUpdateOK  %+v", 200, o.Payload)
}

func (o *MwUpdateOK) GetPayload() *MwUpdateOKBody {
	return o.Payload
}

func (o *MwUpdateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(MwUpdateOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMwUpdateInternalServerError creates a MwUpdateInternalServerError with default headers values
func NewMwUpdateInternalServerError() *MwUpdateInternalServerError {
	return &MwUpdateInternalServerError{}
}

/*
MwUpdateInternalServerError describes a response with status code 500, with default header values.

Generic error specified by `Status` and `Error` fields
*/
type MwUpdateInternalServerError struct {
	Payload *models.Payload
}

// IsSuccess returns true when this mw update internal server error response has a 2xx status code
func (o *MwUpdateInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this mw update internal server error response has a 3xx status code
func (o *MwUpdateInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this mw update internal server error response has a 4xx status code
func (o *MwUpdateInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this mw update internal server error response has a 5xx status code
func (o *MwUpdateInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this mw update internal server error response a status code equal to that given
func (o *MwUpdateInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the mw update internal server error response
func (o *MwUpdateInternalServerError) Code() int {
	return 500
}

func (o *MwUpdateInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /api/mw/{id}][%d] mwUpdateInternalServerError  %+v", 500, o.Payload)
}

func (o *MwUpdateInternalServerError) String() string {
	return fmt.Sprintf("[PUT /api/mw/{id}][%d] mwUpdateInternalServerError  %+v", 500, o.Payload)
}

func (o *MwUpdateInternalServerError) GetPayload() *models.Payload {
	return o.Payload
}

func (o *MwUpdateInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Payload)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
MwUpdateOKBody mw update o k body
swagger:model MwUpdateOKBody
*/
type MwUpdateOKBody struct {

	// error
	Error string `json:"Error,omitempty"`

	// payload
	Payload *MwUpdateOKBodyPayload `json:"Payload,omitempty"`

	// status
	Status string `json:"Status,omitempty"`
}

// Validate validates this mw update o k body
func (o *MwUpdateOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePayload(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *MwUpdateOKBody) validatePayload(formats strfmt.Registry) error {
	if swag.IsZero(o.Payload) { // not required
		return nil
	}

	if o.Payload != nil {
		if err := o.Payload.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mwUpdateOK" + "." + "Payload")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("mwUpdateOK" + "." + "Payload")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this mw update o k body based on the context it is used
func (o *MwUpdateOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidatePayload(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *MwUpdateOKBody) contextValidatePayload(ctx context.Context, formats strfmt.Registry) error {

	if o.Payload != nil {

		if swag.IsZero(o.Payload) { // not required
			return nil
		}

		if err := o.Payload.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mwUpdateOK" + "." + "Payload")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("mwUpdateOK" + "." + "Payload")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *MwUpdateOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *MwUpdateOKBody) UnmarshalBinary(b []byte) error {
	var res MwUpdateOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
MwUpdateOKBodyPayload mw update o k body payload
swagger:model MwUpdateOKBodyPayload
*/
type MwUpdateOKBodyPayload struct {

	// bundle ID
	BundleID string `json:"BundleID,omitempty"`
}

// Validate validates this mw update o k body payload
func (o *MwUpdateOKBodyPayload) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this mw update o k body payload based on context it is used
func (o *MwUpdateOKBodyPayload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *MwUpdateOKBodyPayload) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *MwUpdateOKBodyPayload) UnmarshalBinary(b []byte) error {
	var res MwUpdateOKBodyPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
