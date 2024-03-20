// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetRdmTypeProjectNumberHandlerFunc turns a function with the right signature into a get rdm type project number handler
type GetRdmTypeProjectNumberHandlerFunc func(GetRdmTypeProjectNumberParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetRdmTypeProjectNumberHandlerFunc) Handle(params GetRdmTypeProjectNumberParams) middleware.Responder {
	return fn(params)
}

// GetRdmTypeProjectNumberHandler interface for that can handle valid get rdm type project number params
type GetRdmTypeProjectNumberHandler interface {
	Handle(GetRdmTypeProjectNumberParams) middleware.Responder
}

// NewGetRdmTypeProjectNumber creates a new http.Handler for the get rdm type project number operation
func NewGetRdmTypeProjectNumber(ctx *middleware.Context, handler GetRdmTypeProjectNumberHandler) *GetRdmTypeProjectNumber {
	return &GetRdmTypeProjectNumber{Context: ctx, Handler: handler}
}

/*
	GetRdmTypeProjectNumber swagger:route GET /rdm/{type}/project/{number} getRdmTypeProjectNumber

retrieve RDR data collection associated with a project
*/
type GetRdmTypeProjectNumber struct {
	Context *middleware.Context
	Handler GetRdmTypeProjectNumberHandler
}

func (o *GetRdmTypeProjectNumber) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetRdmTypeProjectNumberParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}