// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// JobData job data
//
// swagger:model jobData
type JobData struct {

	// password of the DR data-access account
	DrPass string `json:"drPass,omitempty"`

	// username of the DR data-access account
	// Required: true
	DrUser *string `json:"drUser"`

	// path or DR namespace (prefixed with irods:) of the destination endpoint
	// Required: true
	DstURL *string `json:"dstURL"`

	// path or DR namespace (prefixed with irods:) of the source endpoint
	// Required: true
	SrcURL *string `json:"srcURL"`

	// username of stager's local account
	// Required: true
	StagerUser *string `json:"stagerUser"`

	// allowed duration in seconds for entire transfer job (0 for no timeout)
	Timeout int64 `json:"timeout,omitempty"`

	// allowed duration in seconds for no further transfer progress (0 for no timeout)
	TimeoutNoprogress int64 `json:"timeout_noprogress,omitempty"`

	// short description about the job
	// Required: true
	Title *string `json:"title"`
}

// Validate validates this job data
func (m *JobData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDrUser(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDstURL(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSrcURL(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStagerUser(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTitle(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JobData) validateDrUser(formats strfmt.Registry) error {

	if err := validate.Required("drUser", "body", m.DrUser); err != nil {
		return err
	}

	return nil
}

func (m *JobData) validateDstURL(formats strfmt.Registry) error {

	if err := validate.Required("dstURL", "body", m.DstURL); err != nil {
		return err
	}

	return nil
}

func (m *JobData) validateSrcURL(formats strfmt.Registry) error {

	if err := validate.Required("srcURL", "body", m.SrcURL); err != nil {
		return err
	}

	return nil
}

func (m *JobData) validateStagerUser(formats strfmt.Registry) error {

	if err := validate.Required("stagerUser", "body", m.StagerUser); err != nil {
		return err
	}

	return nil
}

func (m *JobData) validateTitle(formats strfmt.Registry) error {

	if err := validate.Required("title", "body", m.Title); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this job data based on context it is used
func (m *JobData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *JobData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JobData) UnmarshalBinary(b []byte) error {
	var res JobData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
