// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/dccn-tg/dr-data-stager/pkg/swagger/server/models"
)

// DeleteJobIDOKCode is the HTTP code returned for type DeleteJobIDOK
const DeleteJobIDOKCode int = 200

/*
DeleteJobIDOK success

swagger:response deleteJobIdOK
*/
type DeleteJobIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.JobInfo `json:"body,omitempty"`
}

// NewDeleteJobIDOK creates DeleteJobIDOK with default headers values
func NewDeleteJobIDOK() *DeleteJobIDOK {

	return &DeleteJobIDOK{}
}

// WithPayload adds the payload to the delete job Id o k response
func (o *DeleteJobIDOK) WithPayload(payload *models.JobInfo) *DeleteJobIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete job Id o k response
func (o *DeleteJobIDOK) SetPayload(payload *models.JobInfo) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteJobIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteJobIDNotFoundCode is the HTTP code returned for type DeleteJobIDNotFound
const DeleteJobIDNotFoundCode int = 404

/*
DeleteJobIDNotFound job not found

swagger:response deleteJobIdNotFound
*/
type DeleteJobIDNotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewDeleteJobIDNotFound creates DeleteJobIDNotFound with default headers values
func NewDeleteJobIDNotFound() *DeleteJobIDNotFound {

	return &DeleteJobIDNotFound{}
}

// WithPayload adds the payload to the delete job Id not found response
func (o *DeleteJobIDNotFound) WithPayload(payload string) *DeleteJobIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete job Id not found response
func (o *DeleteJobIDNotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteJobIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// DeleteJobIDInternalServerErrorCode is the HTTP code returned for type DeleteJobIDInternalServerError
const DeleteJobIDInternalServerErrorCode int = 500

/*
DeleteJobIDInternalServerError failure

swagger:response deleteJobIdInternalServerError
*/
type DeleteJobIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ResponseBody500 `json:"body,omitempty"`
}

// NewDeleteJobIDInternalServerError creates DeleteJobIDInternalServerError with default headers values
func NewDeleteJobIDInternalServerError() *DeleteJobIDInternalServerError {

	return &DeleteJobIDInternalServerError{}
}

// WithPayload adds the payload to the delete job Id internal server error response
func (o *DeleteJobIDInternalServerError) WithPayload(payload *models.ResponseBody500) *DeleteJobIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete job Id internal server error response
func (o *DeleteJobIDInternalServerError) SetPayload(payload *models.ResponseBody500) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteJobIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
