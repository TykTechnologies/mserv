// Code generated by go-swagger; DO NOT EDIT.

package mw

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/TykTechnologies/mserv/mservclient/models"
)

// MwFetchReader is a Reader for the MwFetch structure.
type MwFetchReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MwFetchReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewMwFetchOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewMwFetchInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewMwFetchOK creates a MwFetchOK with default headers values
func NewMwFetchOK() *MwFetchOK {
	return &MwFetchOK{}
}

/*MwFetchOK handles this case with default header values.

Full middleware response in the `Payload`
*/
type MwFetchOK struct {
	Payload *MwFetchOKBody
}

func (o *MwFetchOK) Error() string {
	return fmt.Sprintf("[GET /api/mw/{id}][%d] mwFetchOK  %+v", 200, o.Payload)
}

func (o *MwFetchOK) GetPayload() *MwFetchOKBody {
	return o.Payload
}

func (o *MwFetchOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(MwFetchOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMwFetchInternalServerError creates a MwFetchInternalServerError with default headers values
func NewMwFetchInternalServerError() *MwFetchInternalServerError {
	return &MwFetchInternalServerError{}
}

/*MwFetchInternalServerError handles this case with default header values.

Generic error specified by `Status` and `Error` fields
*/
type MwFetchInternalServerError struct {
	Payload *models.Payload
}

func (o *MwFetchInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/mw/{id}][%d] mwFetchInternalServerError  %+v", 500, o.Payload)
}

func (o *MwFetchInternalServerError) GetPayload() *models.Payload {
	return o.Payload
}

func (o *MwFetchInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Payload)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*MwFetchOKBody mw fetch o k body
swagger:model MwFetchOKBody
*/
type MwFetchOKBody struct {

	// error
	Error string `json:"Error,omitempty"`

	// payload
	Payload *models.MW `json:"Payload,omitempty"`

	// status
	Status string `json:"Status,omitempty"`
}

// Validate validates this mw fetch o k body
func (o *MwFetchOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePayload(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *MwFetchOKBody) validatePayload(formats strfmt.Registry) error {

	if swag.IsZero(o.Payload) { // not required
		return nil
	}

	if o.Payload != nil {
		if err := o.Payload.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mwFetchOK" + "." + "Payload")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *MwFetchOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *MwFetchOKBody) UnmarshalBinary(b []byte) error {
	var res MwFetchOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
