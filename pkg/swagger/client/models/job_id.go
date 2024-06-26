// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
)

// JobID identifier for scheduled background tasks.
//
// swagger:model jobID
type JobID string

// Validate validates this job ID
func (m JobID) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this job ID based on context it is used
func (m JobID) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
