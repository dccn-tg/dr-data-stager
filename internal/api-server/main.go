package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
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
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/s12v/go-jwks"
	"github.com/square/go-jose"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
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

	// parse the redis URL to redis connection options
	redisOpts, err := asynq.ParseRedisURI(*redisURL)
	if err != nil {
		log.Fatalf("cannot parse redis URL: %s", err)
	}

	// initialize another redis client for incremental taskId generation
	rdbOpts, _ := redis.ParseURL(*redisURL)
	rdb4tid := redis.NewClient(rdbOpts)
	defer rdb4tid.Close()

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

	client := asynq.NewClient(redisOpts)
	defer client.Close()

	inspector := asynq.NewInspector(redisOpts)

	if err != nil {
		log.Errorf("cannot connect to db: %s", err)
		os.Exit(1)
	}

	// authentication with username/password.
	api.BasicAuthAuth = func(username, password string) (*models.Principal, error) {

		pass, ok := cfg.Auth[username]

		if ok && pass == password {
			principal := models.Principal(username)
			return &principal, nil
		}

		// username not in the static credential list, assuming that the password is a oauth2 access token
		token, err := verifyJwt(password, cfg.Oauth2.JwksEndpoint)
		if err != nil || !token.Valid {
			return nil, errors.New(401, "invalid token: %s", err)
		}

		// there is login user information attached, set the pricipal as the username.
		principal := models.Principal(username)
		return &principal, nil
	}

	// authentication with oauth2 token.
	api.Oauth2Auth = func(tokenStr string, scopes []string) (*models.Principal, error) {

		token, err := verifyJwt(tokenStr, cfg.Oauth2.JwksEndpoint)

		if err != nil || !token.Valid {
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

		uid := claims.ClientID
		// retrieve user profile if "urn:dccn:identity:uid" in the scopes of the token
		// use "urn:dccn:uid" of the profile as the principal.
		if hasScope(claims.Scope, "urn:dccn:identity:uid") {
			uinfo, err := oauth2GetUserInfo(tokenStr, cfg.Oauth2.UserInfoEndpoint)
			if err != nil {
				return nil, errors.New(401, "cannot get userinfo: %s", err)
			}
			uid = uinfo.Uid
		}
		principal := models.Principal(uid)
		return &principal, nil
	}

	api.GetPingHandler = operations.GetPingHandlerFunc(handler.GetPing(cfg))
	api.GetJobIDHandler = operations.GetJobIDHandlerFunc(handler.GetJob(ctx, inspector))
	api.DeleteJobIDHandler = operations.DeleteJobIDHandlerFunc(handler.DeleteJob(ctx, inspector))
	api.PostJobHandler = operations.PostJobHandlerFunc(handler.NewJob(ctx, client, rdb4tid))
	api.PostJobsHandler = operations.PostJobsHandlerFunc(handler.NewJobs(ctx, client, rdb4tid))
	api.GetJobsHandler = operations.GetJobsHandlerFunc(handler.GetJobs(ctx, inspector))
	api.GetDirHandler = operations.GetDirHandlerFunc(handler.ListDir(ctx))
	api.GetCollectionTypeProjectNumberHandler = operations.GetCollectionTypeProjectNumberHandlerFunc(handler.GetCollectionByProject(ctx, cfg))
	api.GetRdmTypeProjectNumberHandler = operations.GetRdmTypeProjectNumberHandlerFunc(handler.GetRdmByProject(ctx, cfg))
	// configure API
	server.ConfigureAPI()

	// Start API server
	if err := server.Serve(); err != nil {
		log.Fatalf("%s", err)
	}
}

func hasScope(scopes []string, scope string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// custom claims data structure, this should match the
// data structure expected from the authentication server.
type IDServerClaims struct {
	Scope    []string `json:"scope"`
	Audience []string `json:"aud"`
	ClientID string   `json:"client_id"`
	jwt.StandardClaims
}

func verifyJwt(tokenStr, jwksEndpoint string) (*jwt.Token, error) {

	return jwt.ParseWithClaims(tokenStr, &IDServerClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(401, "unexpected signing method: %v", token.Header["alg"])
		}

		// get public key from the auth server
		// TODO: discover jwks endpoint using oidc client.
		jwksSource := jwks.NewWebSource(jwksEndpoint)
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
}

type OauthUserInfo struct {
	Uid string `json:"urn:dccn:uid,omitempty"`
}

func oauth2GetUserInfo(tokenStr, url string) (*OauthUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("GET %s (%s) %s\n", url, res.Status, resBody)

	var uinfo OauthUserInfo
	if err := json.Unmarshal(resBody, &uinfo); err != nil {
		return nil, err
	}

	return &uinfo, nil
}
