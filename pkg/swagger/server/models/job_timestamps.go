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

// JobTimestamps job timestamps
//
// swagger:model jobTimestamps
type JobTimestamps struct {

	// timestamp at which the job is completed, -62135596800 (0001-01-01T00:00:00) if not applicable.
	// Required: true
	CompletedAt *int64 `json:"completedAt"`

	// timestamp at which the job is created.
	// Required: true
	CreatedAt *int64 `json:"createdAt"`

	// timestamp at which the job failed the last time, -62135596800 (0001-01-01T00:00:00) if not applicable.
	// Required: true
	LastFailedAt *int64 `json:"lastFailedAt"`

	// timestamp at which the job will be processed, -62135596800 (0001-01-01T00:00:00) if not applicable.
	// Required: true
	NextProcessAt *int64 `json:"nextProcessAt"`
}

// Validate validates this job timestamps
func (m *JobTimestamps) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCompletedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastFailedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNextProcessAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JobTimestamps) validateCompletedAt(formats strfmt.Registry) error {

	if err := validate.Required("completedAt", "body", m.CompletedAt); err != nil {
		return err
	}

	return nil
}

func (m *JobTimestamps) validateCreatedAt(formats strfmt.Registry) error {

	if err := validate.Required("createdAt", "body", m.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (m *JobTimestamps) validateLastFailedAt(formats strfmt.Registry) error {

	if err := validate.Required("lastFailedAt", "body", m.LastFailedAt); err != nil {
		return err
	}

	return nil
}

func (m *JobTimestamps) validateNextProcessAt(formats strfmt.Registry) error {

	if err := validate.Required("nextProcessAt", "body", m.NextProcessAt); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this job timestamps based on context it is used
func (m *JobTimestamps) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *JobTimestamps) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JobTimestamps) UnmarshalBinary(b []byte) error {
	var res JobTimestamps
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
