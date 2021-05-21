// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/client/models"
)

// GetPingReader is a Reader for the GetPing structure.
type GetPingReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPingReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPingOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetPingInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetPingOK creates a GetPingOK with default headers values
func NewGetPingOK() *GetPingOK {
	return &GetPingOK{}
}

/* GetPingOK describes a response with status code 200, with default header values.

success
*/
type GetPingOK struct {
	Payload string
}

func (o *GetPingOK) Error() string {
	return fmt.Sprintf("[GET /ping][%d] getPingOK  %+v", 200, o.Payload)
}
func (o *GetPingOK) GetPayload() string {
	return o.Payload
}

func (o *GetPingOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPingInternalServerError creates a GetPingInternalServerError with default headers values
func NewGetPingInternalServerError() *GetPingInternalServerError {
	return &GetPingInternalServerError{}
}

/* GetPingInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetPingInternalServerError struct {
	Payload *models.ResponseBody500
}

func (o *GetPingInternalServerError) Error() string {
	return fmt.Sprintf("[GET /ping][%d] getPingInternalServerError  %+v", 500, o.Payload)
}
func (o *GetPingInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetPingInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
