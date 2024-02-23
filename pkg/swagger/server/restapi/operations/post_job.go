// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
)

// PostJobHandlerFunc turns a function with the right signature into a post job handler
type PostJobHandlerFunc func(PostJobParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn PostJobHandlerFunc) Handle(params PostJobParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// PostJobHandler interface for that can handle valid post job params
type PostJobHandler interface {
	Handle(PostJobParams, *models.Principal) middleware.Responder
}

// NewPostJob creates a new http.Handler for the post job operation
func NewPostJob(ctx *middleware.Context, handler PostJobHandler) *PostJob {
	return &PostJob{Context: ctx, Handler: handler}
}

/*
	PostJob swagger:route POST /job postJob

create a new stager job
*/
type PostJob struct {
	Context *middleware.Context
	Handler PostJobHandler
}

func (o *PostJob) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostJobParams()
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
