// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/dccn-tg/dr-data-stager/pkg/swagger/server/models"
)

// NewDrDataStagerAPI creates a new DrDataStager instance
func NewDrDataStagerAPI(spec *loads.Document) *DrDataStagerAPI {
	return &DrDataStagerAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		DeleteJobIDHandler: DeleteJobIDHandlerFunc(func(params DeleteJobIDParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation DeleteJobID has not yet been implemented")
		}),
		GetDacProjectNumberHandler: GetDacProjectNumberHandlerFunc(func(params GetDacProjectNumberParams) middleware.Responder {
			return middleware.NotImplemented("operation GetDacProjectNumber has not yet been implemented")
		}),
		GetDirHandler: GetDirHandlerFunc(func(params GetDirParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation GetDir has not yet been implemented")
		}),
		GetJobIDHandler: GetJobIDHandlerFunc(func(params GetJobIDParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation GetJobID has not yet been implemented")
		}),
		GetJobsHandler: GetJobsHandlerFunc(func(params GetJobsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation GetJobs has not yet been implemented")
		}),
		GetJobsStatusHandler: GetJobsStatusHandlerFunc(func(params GetJobsStatusParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation GetJobsStatus has not yet been implemented")
		}),
		GetPingHandler: GetPingHandlerFunc(func(params GetPingParams) middleware.Responder {
			return middleware.NotImplemented("operation GetPing has not yet been implemented")
		}),
		PostJobHandler: PostJobHandlerFunc(func(params PostJobParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation PostJob has not yet been implemented")
		}),
		PostJobsHandler: PostJobsHandlerFunc(func(params PostJobsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation PostJobs has not yet been implemented")
		}),
		PutJobScheduledIDHandler: PutJobScheduledIDHandlerFunc(func(params PutJobScheduledIDParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation PutJobScheduledID has not yet been implemented")
		}),

		// Applies when the Authorization header is set with the Basic scheme
		BasicAuthAuth: func(user string, pass string) (*models.Principal, error) {
			return nil, errors.NotImplemented("basic auth  (basicAuth) has not yet been implemented")
		},
		Oauth2Auth: func(token string, scopes []string) (*models.Principal, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (oauth2) has not yet been implemented")
		},
		// default authorizer is authorized meaning no requests are blocked
		APIAuthorizer: security.Authorized(),
	}
}

/*DrDataStagerAPI Donders Repository data stager APIs */
type DrDataStagerAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// BasicAuthAuth registers a function that takes username and password and returns a principal
	// it performs authentication with basic auth
	BasicAuthAuth func(string, string) (*models.Principal, error)

	// Oauth2Auth registers a function that takes an access token and a collection of required scopes and returns a principal
	// it performs authentication based on an oauth2 bearer token provided in the request
	Oauth2Auth func(string, []string) (*models.Principal, error)

	// APIAuthorizer provides access control (ACL/RBAC/ABAC) by providing access to the request and authenticated principal
	APIAuthorizer runtime.Authorizer

	// DeleteJobIDHandler sets the operation handler for the delete job ID operation
	DeleteJobIDHandler DeleteJobIDHandler
	// GetDacProjectNumberHandler sets the operation handler for the get dac project number operation
	GetDacProjectNumberHandler GetDacProjectNumberHandler
	// GetDirHandler sets the operation handler for the get dir operation
	GetDirHandler GetDirHandler
	// GetJobIDHandler sets the operation handler for the get job ID operation
	GetJobIDHandler GetJobIDHandler
	// GetJobsHandler sets the operation handler for the get jobs operation
	GetJobsHandler GetJobsHandler
	// GetJobsStatusHandler sets the operation handler for the get jobs status operation
	GetJobsStatusHandler GetJobsStatusHandler
	// GetPingHandler sets the operation handler for the get ping operation
	GetPingHandler GetPingHandler
	// PostJobHandler sets the operation handler for the post job operation
	PostJobHandler PostJobHandler
	// PostJobsHandler sets the operation handler for the post jobs operation
	PostJobsHandler PostJobsHandler
	// PutJobScheduledIDHandler sets the operation handler for the put job scheduled ID operation
	PutJobScheduledIDHandler PutJobScheduledIDHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *DrDataStagerAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *DrDataStagerAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *DrDataStagerAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *DrDataStagerAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *DrDataStagerAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *DrDataStagerAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *DrDataStagerAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *DrDataStagerAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *DrDataStagerAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the DrDataStagerAPI
func (o *DrDataStagerAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.BasicAuthAuth == nil {
		unregistered = append(unregistered, "BasicAuthAuth")
	}
	if o.Oauth2Auth == nil {
		unregistered = append(unregistered, "Oauth2Auth")
	}

	if o.DeleteJobIDHandler == nil {
		unregistered = append(unregistered, "DeleteJobIDHandler")
	}
	if o.GetDacProjectNumberHandler == nil {
		unregistered = append(unregistered, "GetDacProjectNumberHandler")
	}
	if o.GetDirHandler == nil {
		unregistered = append(unregistered, "GetDirHandler")
	}
	if o.GetJobIDHandler == nil {
		unregistered = append(unregistered, "GetJobIDHandler")
	}
	if o.GetJobsHandler == nil {
		unregistered = append(unregistered, "GetJobsHandler")
	}
	if o.GetJobsStatusHandler == nil {
		unregistered = append(unregistered, "GetJobsStatusHandler")
	}
	if o.GetPingHandler == nil {
		unregistered = append(unregistered, "GetPingHandler")
	}
	if o.PostJobHandler == nil {
		unregistered = append(unregistered, "PostJobHandler")
	}
	if o.PostJobsHandler == nil {
		unregistered = append(unregistered, "PostJobsHandler")
	}
	if o.PutJobScheduledIDHandler == nil {
		unregistered = append(unregistered, "PutJobScheduledIDHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *DrDataStagerAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *DrDataStagerAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	result := make(map[string]runtime.Authenticator)
	for name := range schemes {
		switch name {
		case "basicAuth":
			result[name] = o.BasicAuthenticator(func(username, password string) (interface{}, error) {
				return o.BasicAuthAuth(username, password)
			})

		case "oauth2":
			result[name] = o.BearerAuthenticator(name, func(token string, scopes []string) (interface{}, error) {
				return o.Oauth2Auth(token, scopes)
			})

		}
	}
	return result
}

// Authorizer returns the registered authorizer
func (o *DrDataStagerAPI) Authorizer() runtime.Authorizer {
	return o.APIAuthorizer
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *DrDataStagerAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *DrDataStagerAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *DrDataStagerAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the dr data stager API
func (o *DrDataStagerAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *DrDataStagerAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/job/{id}"] = NewDeleteJobID(o.context, o.DeleteJobIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/dac/project/{number}"] = NewGetDacProjectNumber(o.context, o.GetDacProjectNumberHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/dir"] = NewGetDir(o.context, o.GetDirHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/job/{id}"] = NewGetJobID(o.context, o.GetJobIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/jobs"] = NewGetJobs(o.context, o.GetJobsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/jobs/{status}"] = NewGetJobsStatus(o.context, o.GetJobsStatusHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/ping"] = NewGetPing(o.context, o.GetPingHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/job"] = NewPostJob(o.context, o.PostJobHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/jobs"] = NewPostJobs(o.context, o.PostJobsHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/job/scheduled/{id}"] = NewPutJobScheduledID(o.context, o.PutJobScheduledIDHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *DrDataStagerAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *DrDataStagerAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *DrDataStagerAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *DrDataStagerAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *DrDataStagerAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[um][path] = builder(h)
	}
}
