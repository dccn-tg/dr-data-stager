// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DirEntry directory entry
//
// swagger:model dirEntry
type DirEntry struct {

	// name of the entry
	// Required: true
	Name *string `json:"name"`

	// size of the entry in bytes
	// Required: true
	Size *int64 `json:"size"`

	// type of the entry
	// Required: true
	// Enum: [regular dir symlink unknown]
	Type *string `json:"type"`
}

// Validate validates this dir entry
func (m *DirEntry) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSize(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DirEntry) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *DirEntry) validateSize(formats strfmt.Registry) error {

	if err := validate.Required("size", "body", m.Size); err != nil {
		return err
	}

	return nil
}

var dirEntryTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["regular","dir","symlink","unknown"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		dirEntryTypeTypePropEnum = append(dirEntryTypeTypePropEnum, v)
	}
}

const (

	// DirEntryTypeRegular captures enum value "regular"
	DirEntryTypeRegular string = "regular"

	// DirEntryTypeDir captures enum value "dir"
	DirEntryTypeDir string = "dir"

	// DirEntryTypeSymlink captures enum value "symlink"
	DirEntryTypeSymlink string = "symlink"

	// DirEntryTypeUnknown captures enum value "unknown"
	DirEntryTypeUnknown string = "unknown"
)

// prop value enum
func (m *DirEntry) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, dirEntryTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *DirEntry) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this dir entry based on context it is used
func (m *DirEntry) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DirEntry) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DirEntry) UnmarshalBinary(b []byte) error {
	var res DirEntry
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}