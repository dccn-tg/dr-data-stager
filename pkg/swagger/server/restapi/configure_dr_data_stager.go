// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name DrDataStager --spec ../../swagger.yaml --principal models.Principal --exclude-main

func configureFlags(api *operations.DrDataStagerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.DrDataStagerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the Authorization header is set with the Basic scheme
	if api.BasicAuthAuth == nil {
		api.BasicAuthAuth = func(user string, pass string) (*models.Principal, error) {
			return nil, errors.NotImplemented("basic auth  (basicAuth) has not yet been implemented")
		}
	}
	if api.Oauth2Auth == nil {
		api.Oauth2Auth = func(token string, scopes []string) (*models.Principal, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (oauth2) has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.DeleteJobIDHandler == nil {
		api.DeleteJobIDHandler = operations.DeleteJobIDHandlerFunc(func(params operations.DeleteJobIDParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteJobID has not yet been implemented")
		})
	}
	if api.GetDacProjectNumberHandler == nil {
		api.GetDacProjectNumberHandler = operations.GetDacProjectNumberHandlerFunc(func(params operations.GetDacProjectNumberParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetDacProjectNumber has not yet been implemented")
		})
	}
	if api.GetDirHandler == nil {
		api.GetDirHandler = operations.GetDirHandlerFunc(func(params operations.GetDirParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetDir has not yet been implemented")
		})
	}
	if api.GetJobIDHandler == nil {
		api.GetJobIDHandler = operations.GetJobIDHandlerFunc(func(params operations.GetJobIDParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetJobID has not yet been implemented")
		})
	}
	if api.GetJobsHandler == nil {
		api.GetJobsHandler = operations.GetJobsHandlerFunc(func(params operations.GetJobsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetJobs has not yet been implemented")
		})
	}
	if api.GetJobsStatusHandler == nil {
		api.GetJobsStatusHandler = operations.GetJobsStatusHandlerFunc(func(params operations.GetJobsStatusParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetJobsStatus has not yet been implemented")
		})
	}
	if api.GetPingHandler == nil {
		api.GetPingHandler = operations.GetPingHandlerFunc(func(params operations.GetPingParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPing has not yet been implemented")
		})
	}
	if api.PostJobHandler == nil {
		api.PostJobHandler = operations.PostJobHandlerFunc(func(params operations.PostJobParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostJob has not yet been implemented")
		})
	}
	if api.PostJobsHandler == nil {
		api.PostJobsHandler = operations.PostJobsHandlerFunc(func(params operations.PostJobsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostJobs has not yet been implemented")
		})
	}
	if api.PutJobScheduledIDHandler == nil {
		api.PutJobScheduledIDHandler = operations.PutJobScheduledIDHandlerFunc(func(params operations.PutJobScheduledIDParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation operations.PutJobScheduledID has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
