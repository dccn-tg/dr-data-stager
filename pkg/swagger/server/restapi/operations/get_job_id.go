// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
)

// GetJobIDHandlerFunc turns a function with the right signature into a get job ID handler
type GetJobIDHandlerFunc func(GetJobIDParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetJobIDHandlerFunc) Handle(params GetJobIDParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetJobIDHandler interface for that can handle valid get job ID params
type GetJobIDHandler interface {
	Handle(GetJobIDParams, *models.Principal) middleware.Responder
}

// NewGetJobID creates a new http.Handler for the get job ID operation
func NewGetJobID(ctx *middleware.Context, handler GetJobIDHandler) *GetJobID {
	return &GetJobID{Context: ctx, Handler: handler}
}

/*
	GetJobID swagger:route GET /job/{id} getJobId

get stager job information
*/
type GetJobID struct {
	Context *middleware.Context
	Handler GetJobIDHandler
}

func (o *GetJobID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetJobIDParams()
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
