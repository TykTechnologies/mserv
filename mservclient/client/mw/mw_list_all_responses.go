// Code generated by go-swagger; DO NOT EDIT.

package mw

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/TykTechnologies/mserv/mservclient/models"
)

// MwListAllReader is a Reader for the MwListAll structure.
type MwListAllReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MwListAllReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewMwListAllOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewMwListAllInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /api/mw/master/all] mwListAll", response, response.Code())
	}
}

// NewMwListAllOK creates a MwListAllOK with default headers values
func NewMwListAllOK() *MwListAllOK {
	return &MwListAllOK{}
}

/*
MwListAllOK describes a response with status code 200, with default header values.

List of full middleware representations in the `Payload`
*/
type MwListAllOK struct {
	Payload *MwListAllOKBody
}

// IsSuccess returns true when this mw list all o k response has a 2xx status code
func (o *MwListAllOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this mw list all o k response has a 3xx status code
func (o *MwListAllOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this mw list all o k response has a 4xx status code
func (o *MwListAllOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this mw list all o k response has a 5xx status code
func (o *MwListAllOK) IsServerError() bool {
	return false
}

// IsCode returns true when this mw list all o k response a status code equal to that given
func (o *MwListAllOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the mw list all o k response
func (o *MwListAllOK) Code() int {
	return 200
}

func (o *MwListAllOK) Error() string {
	return fmt.Sprintf("[GET /api/mw/master/all][%d] mwListAllOK  %+v", 200, o.Payload)
}

func (o *MwListAllOK) String() string {
	return fmt.Sprintf("[GET /api/mw/master/all][%d] mwListAllOK  %+v", 200, o.Payload)
}

func (o *MwListAllOK) GetPayload() *MwListAllOKBody {
	return o.Payload
}

func (o *MwListAllOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(MwListAllOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMwListAllInternalServerError creates a MwListAllInternalServerError with default headers values
func NewMwListAllInternalServerError() *MwListAllInternalServerError {
	return &MwListAllInternalServerError{}
}

/*
MwListAllInternalServerError describes a response with status code 500, with default header values.

Generic error specified by `Status` and `Error` fields
*/
type MwListAllInternalServerError struct {
	Payload *models.Payload
}

// IsSuccess returns true when this mw list all internal server error response has a 2xx status code
func (o *MwListAllInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this mw list all internal server error response has a 3xx status code
func (o *MwListAllInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this mw list all internal server error response has a 4xx status code
func (o *MwListAllInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this mw list all internal server error response has a 5xx status code
func (o *MwListAllInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this mw list all internal server error response a status code equal to that given
func (o *MwListAllInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the mw list all internal server error response
func (o *MwListAllInternalServerError) Code() int {
	return 500
}

func (o *MwListAllInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/mw/master/all][%d] mwListAllInternalServerError  %+v", 500, o.Payload)
}

func (o *MwListAllInternalServerError) String() string {
	return fmt.Sprintf("[GET /api/mw/master/all][%d] mwListAllInternalServerError  %+v", 500, o.Payload)
}

func (o *MwListAllInternalServerError) GetPayload() *models.Payload {
	return o.Payload
}

func (o *MwListAllInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Payload)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
MwListAllOKBody mw list all o k body
swagger:model MwListAllOKBody
*/
type MwListAllOKBody struct {

	// error
	Error string `json:"Error,omitempty"`

	// payload
	Payload []*models.MW `json:"Payload"`

	// status
	Status string `json:"Status,omitempty"`
}

// Validate validates this mw list all o k body
func (o *MwListAllOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePayload(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *MwListAllOKBody) validatePayload(formats strfmt.Registry) error {
	if swag.IsZero(o.Payload) { // not required
		return nil
	}

	for i := 0; i < len(o.Payload); i++ {
		if swag.IsZero(o.Payload[i]) { // not required
			continue
		}

		if o.Payload[i] != nil {
			if err := o.Payload[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("mwListAllOK" + "." + "Payload" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("mwListAllOK" + "." + "Payload" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this mw list all o k body based on the context it is used
func (o *MwListAllOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidatePayload(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *MwListAllOKBody) contextValidatePayload(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Payload); i++ {

		if o.Payload[i] != nil {

			if swag.IsZero(o.Payload[i]) { // not required
				return nil
			}

			if err := o.Payload[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("mwListAllOK" + "." + "Payload" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("mwListAllOK" + "." + "Payload" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *MwListAllOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *MwListAllOKBody) UnmarshalBinary(b []byte) error {
	var res MwListAllOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
