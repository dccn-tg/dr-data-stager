package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Donders-Institute/dr-data-stager/internal/api-server/config"
	"github.com/Donders-Institute/dr-data-stager/internal/api-server/handler"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/restapi"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/restapi/operations"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-redis/redis/v8"
	"github.com/s12v/go-jwks"
	"github.com/square/go-jose"

	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/logging"
	"github.com/thoas/bokchoy/middleware"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	optsPort    *int
	redisURL    *string
	configFile  *string
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	optsPort = flag.Int("p", 8080, "specify the service `port` number")
	redisURL = flag.String("r", "redis://redis:6379", "redis service `address`")
	configFile = flag.String("c", os.Getenv("STAGER_APISERVER_CONFIG"), "configurateion file `path`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/api-server.log",
		FileLevel:         log.Info,
	}

	if *optsVerbose {
		cfg.ConsoleLevel = log.Debug
		cfg.FileLevel = log.Debug
	}

	// initialize logger
	log.NewLogger(cfg, log.InstanceZapLogger)
}

func usage() {
	fmt.Printf("\nAPI server of filer gateway\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	// load global configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("fail to load configuration: %s", *configFile)
	}

	// redis client instance for notifying cache update
	redisOpts, err := redis.ParseURL(*redisURL)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalf("%s", err)
	}

	api := operations.NewDrDataStagerAPI(swaggerSpec)
	api.UseRedoc()
	server := restapi.NewServer(api)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// actions to take when the main program exists.
	defer func() {

		// stop all background context
		cancel()

		// stop API server.
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalf("%s", err)
		}
	}()

	server.Port = *optsPort
	server.ListenLimit = 10
	server.TLSListenLimit = 10

	// initiate blochy queue for setting project roles
	var logger logging.Logger

	bok, err := bokchoy.New(ctx, bokchoy.Config{
		Broker: bokchoy.BrokerConfig{
			Type: "redis",
			Redis: bokchoy.RedisConfig{
				Type: "client",
				Client: bokchoy.RedisClientConfig{
					Addr: redisOpts.Addr,
				},
			},
		},
	}, bokchoy.WithLogger(logger), bokchoy.WithTTL(7*24*time.Hour))

	if err != nil {
		log.Errorf("cannot connect to db: %s", err)
		os.Exit(1)
	}

	bok.Use(middleware.Recoverer)
	bok.Use(middleware.DefaultLogger)

	// authentication with username/password.
	api.BasicAuthAuth = func(username, password string) (*models.Principal, error) {

		pass, ok := cfg.Auth[username]

		if !ok || pass != password {
			return nil, errors.New(401, "incorrect username/password")
		}

		// there is login user information attached, set the pricipal as the username.
		principal := models.Principal(username)
		return &principal, nil
	}

	// authentication with oauth2 token.
	api.Oauth2Auth = func(tokenStr string, scopes []string) (*models.Principal, error) {

		// custom claims data structure, this should match the
		// data structure expected from the authentication server.
		type IDServerClaims struct {
			Scope    []string `json:"scope"`
			Audience []string `json:"aud"`
			ClientID string   `json:"client_id"`
			jwt.StandardClaims
		}

		token, err := jwt.ParseWithClaims(tokenStr, &IDServerClaims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New(401, "unexpected signing method: %v", token.Header["alg"])
			}

			// get public key from the auth server
			// TODO: discover jwks endpoint using oidc client.
			jwksSource := jwks.NewWebSource(cfg.JwksEndpoint)
			jwksClient := jwks.NewDefaultClient(
				jwksSource,
				time.Hour,    // Refresh keys every 1 hour
				12*time.Hour, // Expire keys after 12 hours
			)

			var jwk *jose.JSONWebKey
			jwk, err := jwksClient.GetEncryptionKey(token.Header["kid"].(string))
			if err != nil {
				return nil, errors.New(401, "cannot retrieve encryption key: %s", err)
			}

			return jwk.Key, nil
		})

		if err != nil {
			return nil, errors.New(401, "invalid token: %s", err)
		}

		// check token scope
		claims, ok := token.Claims.(*IDServerClaims)
		if !ok {
			return nil, errors.New(401, "cannot get claims from the token")
		}

		inScope := func(target string) bool {
			for _, s := range claims.Scope {
				if s == target {
					return true
				}
			}
			return false
		}

		for _, scope := range scopes {
			if !inScope(scope) {
				return nil, errors.New(401, "token not in scope: %s", scope)
			}
		}

		principal := models.Principal(claims.ClientID)
		return &principal, nil
	}

	api.GetPingHandler = operations.GetPingHandlerFunc(handler.GetPing(cfg))
	api.GetJobIDHandler = operations.GetJobIDHandlerFunc(handler.GetJob(ctx, bok))
	api.PostJobHandler = operations.PostJobHandlerFunc(handler.NewJob(ctx, bok))

	// configure API
	server.ConfigureAPI()

	// Start API server
	if err := server.Serve(); err != nil {
		log.Fatalf("%s", err)
	}
}
