// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewDeleteJobIDParams creates a new DeleteJobIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteJobIDParams() *DeleteJobIDParams {
	return &DeleteJobIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteJobIDParamsWithTimeout creates a new DeleteJobIDParams object
// with the ability to set a timeout on a request.
func NewDeleteJobIDParamsWithTimeout(timeout time.Duration) *DeleteJobIDParams {
	return &DeleteJobIDParams{
		timeout: timeout,
	}
}

// NewDeleteJobIDParamsWithContext creates a new DeleteJobIDParams object
// with the ability to set a context for a request.
func NewDeleteJobIDParamsWithContext(ctx context.Context) *DeleteJobIDParams {
	return &DeleteJobIDParams{
		Context: ctx,
	}
}

// NewDeleteJobIDParamsWithHTTPClient creates a new DeleteJobIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteJobIDParamsWithHTTPClient(client *http.Client) *DeleteJobIDParams {
	return &DeleteJobIDParams{
		HTTPClient: client,
	}
}

/*
DeleteJobIDParams contains all the parameters to send to the API endpoint

	for the delete job ID operation.

	Typically these are written to a http.Request.
*/
type DeleteJobIDParams struct {

	/* ID.

	   job identifier
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete job ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteJobIDParams) WithDefaults() *DeleteJobIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete job ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteJobIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete job ID params
func (o *DeleteJobIDParams) WithTimeout(timeout time.Duration) *DeleteJobIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete job ID params
func (o *DeleteJobIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete job ID params
func (o *DeleteJobIDParams) WithContext(ctx context.Context) *DeleteJobIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete job ID params
func (o *DeleteJobIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete job ID params
func (o *DeleteJobIDParams) WithHTTPClient(client *http.Client) *DeleteJobIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete job ID params
func (o *DeleteJobIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the delete job ID params
func (o *DeleteJobIDParams) WithID(id string) *DeleteJobIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete job ID params
func (o *DeleteJobIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteJobIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
