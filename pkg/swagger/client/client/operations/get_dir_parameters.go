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

// NewGetDirParams creates a new GetDirParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetDirParams() *GetDirParams {
	return &GetDirParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetDirParamsWithTimeout creates a new GetDirParams object
// with the ability to set a timeout on a request.
func NewGetDirParamsWithTimeout(timeout time.Duration) *GetDirParams {
	return &GetDirParams{
		timeout: timeout,
	}
}

// NewGetDirParamsWithContext creates a new GetDirParams object
// with the ability to set a context for a request.
func NewGetDirParamsWithContext(ctx context.Context) *GetDirParams {
	return &GetDirParams{
		Context: ctx,
	}
}

// NewGetDirParamsWithHTTPClient creates a new GetDirParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetDirParamsWithHTTPClient(client *http.Client) *GetDirParams {
	return &GetDirParams{
		HTTPClient: client,
	}
}

/*
GetDirParams contains all the parameters to send to the API endpoint

	for the get dir operation.

	Typically these are written to a http.Request.
*/
type GetDirParams struct {

	/* Path.

	   path
	*/
	Path string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get dir params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDirParams) WithDefaults() *GetDirParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get dir params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDirParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get dir params
func (o *GetDirParams) WithTimeout(timeout time.Duration) *GetDirParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get dir params
func (o *GetDirParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get dir params
func (o *GetDirParams) WithContext(ctx context.Context) *GetDirParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get dir params
func (o *GetDirParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get dir params
func (o *GetDirParams) WithHTTPClient(client *http.Client) *GetDirParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get dir params
func (o *GetDirParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPath adds the path to the get dir params
func (o *GetDirParams) WithPath(path string) *GetDirParams {
	o.SetPath(path)
	return o
}

// SetPath adds the path to the get dir params
func (o *GetDirParams) SetPath(path string) {
	o.Path = path
}

// WriteToRequest writes these params to a swagger request
func (o *GetDirParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param path
	qrPath := o.Path
	qPath := qrPath
	if qPath != "" {

		if err := r.SetQueryParam("path", qPath); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
