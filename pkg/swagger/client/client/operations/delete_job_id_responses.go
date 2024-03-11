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

// DeleteJobIDReader is a Reader for the DeleteJobID structure.
type DeleteJobIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteJobIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteJobIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteJobIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDeleteJobIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteJobIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[DELETE /job/{id}] DeleteJobID", response, response.Code())
	}
}

// NewDeleteJobIDOK creates a DeleteJobIDOK with default headers values
func NewDeleteJobIDOK() *DeleteJobIDOK {
	return &DeleteJobIDOK{}
}

/*
DeleteJobIDOK describes a response with status code 200, with default header values.

success
*/
type DeleteJobIDOK struct {
	Payload *models.JobInfo
}

// IsSuccess returns true when this delete job Id o k response has a 2xx status code
func (o *DeleteJobIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete job Id o k response has a 3xx status code
func (o *DeleteJobIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete job Id o k response has a 4xx status code
func (o *DeleteJobIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete job Id o k response has a 5xx status code
func (o *DeleteJobIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete job Id o k response a status code equal to that given
func (o *DeleteJobIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete job Id o k response
func (o *DeleteJobIDOK) Code() int {
	return 200
}

func (o *DeleteJobIDOK) Error() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdOK  %+v", 200, o.Payload)
}

func (o *DeleteJobIDOK) String() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdOK  %+v", 200, o.Payload)
}

func (o *DeleteJobIDOK) GetPayload() *models.JobInfo {
	return o.Payload
}

func (o *DeleteJobIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.JobInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteJobIDBadRequest creates a DeleteJobIDBadRequest with default headers values
func NewDeleteJobIDBadRequest() *DeleteJobIDBadRequest {
	return &DeleteJobIDBadRequest{}
}

/*
DeleteJobIDBadRequest describes a response with status code 400, with default header values.

bad request
*/
type DeleteJobIDBadRequest struct {
	Payload *models.ResponseBody400
}

// IsSuccess returns true when this delete job Id bad request response has a 2xx status code
func (o *DeleteJobIDBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete job Id bad request response has a 3xx status code
func (o *DeleteJobIDBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete job Id bad request response has a 4xx status code
func (o *DeleteJobIDBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete job Id bad request response has a 5xx status code
func (o *DeleteJobIDBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this delete job Id bad request response a status code equal to that given
func (o *DeleteJobIDBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the delete job Id bad request response
func (o *DeleteJobIDBadRequest) Code() int {
	return 400
}

func (o *DeleteJobIDBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdBadRequest  %+v", 400, o.Payload)
}

func (o *DeleteJobIDBadRequest) String() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdBadRequest  %+v", 400, o.Payload)
}

func (o *DeleteJobIDBadRequest) GetPayload() *models.ResponseBody400 {
	return o.Payload
}

func (o *DeleteJobIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody400)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteJobIDNotFound creates a DeleteJobIDNotFound with default headers values
func NewDeleteJobIDNotFound() *DeleteJobIDNotFound {
	return &DeleteJobIDNotFound{}
}

/*
DeleteJobIDNotFound describes a response with status code 404, with default header values.

job not found
*/
type DeleteJobIDNotFound struct {
	Payload string
}

// IsSuccess returns true when this delete job Id not found response has a 2xx status code
func (o *DeleteJobIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete job Id not found response has a 3xx status code
func (o *DeleteJobIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete job Id not found response has a 4xx status code
func (o *DeleteJobIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete job Id not found response has a 5xx status code
func (o *DeleteJobIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete job Id not found response a status code equal to that given
func (o *DeleteJobIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete job Id not found response
func (o *DeleteJobIDNotFound) Code() int {
	return 404
}

func (o *DeleteJobIDNotFound) Error() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteJobIDNotFound) String() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteJobIDNotFound) GetPayload() string {
	return o.Payload
}

func (o *DeleteJobIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteJobIDInternalServerError creates a DeleteJobIDInternalServerError with default headers values
func NewDeleteJobIDInternalServerError() *DeleteJobIDInternalServerError {
	return &DeleteJobIDInternalServerError{}
}

/*
DeleteJobIDInternalServerError describes a response with status code 500, with default header values.

failure
*/
type DeleteJobIDInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this delete job Id internal server error response has a 2xx status code
func (o *DeleteJobIDInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete job Id internal server error response has a 3xx status code
func (o *DeleteJobIDInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete job Id internal server error response has a 4xx status code
func (o *DeleteJobIDInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete job Id internal server error response has a 5xx status code
func (o *DeleteJobIDInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this delete job Id internal server error response a status code equal to that given
func (o *DeleteJobIDInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the delete job Id internal server error response
func (o *DeleteJobIDInternalServerError) Code() int {
	return 500
}

func (o *DeleteJobIDInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdInternalServerError  %+v", 500, o.Payload)
}

func (o *DeleteJobIDInternalServerError) String() string {
	return fmt.Sprintf("[DELETE /job/{id}][%d] deleteJobIdInternalServerError  %+v", 500, o.Payload)
}

func (o *DeleteJobIDInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *DeleteJobIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
