// Code generated by go-swagger; DO NOT EDIT.

package mw

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/TykTechnologies/mserv/mservclient/models"
)

// MwFetchBundleReader is a Reader for the MwFetchBundle structure.
type MwFetchBundleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MwFetchBundleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewMwFetchBundleOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewMwFetchBundleInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewMwFetchBundleOK creates a MwFetchBundleOK with default headers values
func NewMwFetchBundleOK() *MwFetchBundleOK {
	return &MwFetchBundleOK{}
}

/*MwFetchBundleOK handles this case with default header values.

Middleware bundle as a file
*/
type MwFetchBundleOK struct {
	Payload *models.File
}

func (o *MwFetchBundleOK) Error() string {
	return fmt.Sprintf("[GET /api/mw/bundle/{id}][%d] mwFetchBundleOK  %+v", 200, o.Payload)
}

func (o *MwFetchBundleOK) GetPayload() *models.File {
	return o.Payload
}

func (o *MwFetchBundleOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.File)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMwFetchBundleInternalServerError creates a MwFetchBundleInternalServerError with default headers values
func NewMwFetchBundleInternalServerError() *MwFetchBundleInternalServerError {
	return &MwFetchBundleInternalServerError{}
}

/*MwFetchBundleInternalServerError handles this case with default header values.

Generic error specified by `Status` and `Error` fields
*/
type MwFetchBundleInternalServerError struct {
	Payload *models.Payload
}

func (o *MwFetchBundleInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/mw/bundle/{id}][%d] mwFetchBundleInternalServerError  %+v", 500, o.Payload)
}

func (o *MwFetchBundleInternalServerError) GetPayload() *models.Payload {
	return o.Payload
}

func (o *MwFetchBundleInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Payload)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}