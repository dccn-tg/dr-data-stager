// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
)

// GetDirHandlerFunc turns a function with the right signature into a get dir handler
type GetDirHandlerFunc func(GetDirParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetDirHandlerFunc) Handle(params GetDirParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetDirHandler interface for that can handle valid get dir params
type GetDirHandler interface {
	Handle(GetDirParams, *models.Principal) middleware.Responder
}

// NewGetDir creates a new http.Handler for the get dir operation
func NewGetDir(ctx *middleware.Context, handler GetDirHandler) *GetDir {
	return &GetDir{Context: ctx, Handler: handler}
}

/*
	GetDir swagger:route GET /dir getDir

get entities within a filesystem path
*/
type GetDir struct {
	Context *middleware.Context
	Handler GetDirHandler
}

func (o *GetDir) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetDirParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
