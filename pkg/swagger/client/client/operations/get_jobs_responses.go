// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/dccn-tg/dr-data-stager/pkg/swagger/client/models"
)

// GetJobsReader is a Reader for the GetJobs structure.
type GetJobsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetJobsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetJobsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetJobsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /jobs] GetJobs", response, response.Code())
	}
}

// NewGetJobsOK creates a GetJobsOK with default headers values
func NewGetJobsOK() *GetJobsOK {
	return &GetJobsOK{}
}

/*
GetJobsOK describes a response with status code 200, with default header values.

success
*/
type GetJobsOK struct {
	Payload *models.ResponseBodyJobs
}

// IsSuccess returns true when this get jobs o k response has a 2xx status code
func (o *GetJobsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get jobs o k response has a 3xx status code
func (o *GetJobsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get jobs o k response has a 4xx status code
func (o *GetJobsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get jobs o k response has a 5xx status code
func (o *GetJobsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get jobs o k response a status code equal to that given
func (o *GetJobsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get jobs o k response
func (o *GetJobsOK) Code() int {
	return 200
}

func (o *GetJobsOK) Error() string {
	return fmt.Sprintf("[GET /jobs][%d] getJobsOK  %+v", 200, o.Payload)
}

func (o *GetJobsOK) String() string {
	return fmt.Sprintf("[GET /jobs][%d] getJobsOK  %+v", 200, o.Payload)
}

func (o *GetJobsOK) GetPayload() *models.ResponseBodyJobs {
	return o.Payload
}

func (o *GetJobsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyJobs)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetJobsInternalServerError creates a GetJobsInternalServerError with default headers values
func NewGetJobsInternalServerError() *GetJobsInternalServerError {
	return &GetJobsInternalServerError{}
}

/*
GetJobsInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetJobsInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this get jobs internal server error response has a 2xx status code
func (o *GetJobsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get jobs internal server error response has a 3xx status code
func (o *GetJobsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get jobs internal server error response has a 4xx status code
func (o *GetJobsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get jobs internal server error response has a 5xx status code
func (o *GetJobsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get jobs internal server error response a status code equal to that given
func (o *GetJobsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get jobs internal server error response
func (o *GetJobsInternalServerError) Code() int {
	return 500
}

func (o *GetJobsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /jobs][%d] getJobsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetJobsInternalServerError) String() string {
	return fmt.Sprintf("[GET /jobs][%d] getJobsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetJobsInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetJobsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
