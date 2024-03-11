// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
)

// GetJobsHandlerFunc turns a function with the right signature into a get jobs handler
type GetJobsHandlerFunc func(GetJobsParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetJobsHandlerFunc) Handle(params GetJobsParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetJobsHandler interface for that can handle valid get jobs params
type GetJobsHandler interface {
	Handle(GetJobsParams, *models.Principal) middleware.Responder
}

// NewGetJobs creates a new http.Handler for the get jobs operation
func NewGetJobs(ctx *middleware.Context, handler GetJobsHandler) *GetJobs {
	return &GetJobs{Context: ctx, Handler: handler}
}

/*
	GetJobs swagger:route GET /jobs getJobs

get all jobs of a user
*/
type GetJobs struct {
	Context *middleware.Context
	Handler GetJobsHandler
}

func (o *GetJobs) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetJobsParams()
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
