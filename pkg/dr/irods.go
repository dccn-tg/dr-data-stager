package dr

import (
	"strings"

	"github.com/cyverse/go-irodsclient/fs"
	"github.com/cyverse/go-irodsclient/irods/types"
)

type Config struct {
	IrodsHost           string
	IrodsPort           int
	IrodsZone           string
	IrodsUser           string
	IrodsPass           string
	IrodsAuthScheme     string
	IrodsSslCacert      string
	IrodsSslKeysize     int
	IrodsSslAlgorithm   string
	IrodsSslSaltSize    int
	IrodsHashRounds     int
	OrganisationalUnits []string
}

func (c Config) AuthSchemeType() types.AuthScheme {
	switch strings.ToLower(c.IrodsAuthScheme) {
	case "pam":
		return types.AuthSchemePAM
	case "gsi":
		return types.AuthSchemeGSI
	default:
		return types.AuthSchemeNative
	}
}

func NewAccount(config Config) (*types.IRODSAccount, error) {

	account, err := types.CreateIRODSAccount(
		config.IrodsHost,
		config.IrodsPort,
		config.IrodsUser,
		config.IrodsZone,
		config.AuthSchemeType(),
		config.IrodsPass,
		"",
	)

	if err != nil {
		return nil, err
	}

	sslConfig, err := types.CreateIRODSSSLConfig(
		config.IrodsSslCacert,
		config.IrodsSslKeysize,
		config.IrodsSslAlgorithm,
		config.IrodsSslSaltSize,
		config.IrodsHashRounds,
	)
	if err != nil {
		return nil, err
	}

	account.SSLConfiguration = sslConfig
	account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
	account.ClientServerNegotiation = true

	return account, nil
}

func NewProxyAccount(config Config, user string) (*types.IRODSAccount, error) {

	account, err := types.CreateIRODSProxyAccount(
		config.IrodsHost,
		config.IrodsPort,
		user,
		config.IrodsZone,
		config.IrodsUser,
		config.IrodsZone,
		config.AuthSchemeType(),
		config.IrodsPass,
		"",
	)

	if err != nil {
		return nil, err
	}

	sslConfig, err := types.CreateIRODSSSLConfig(
		config.IrodsSslCacert,
		config.IrodsSslKeysize,
		config.IrodsSslAlgorithm,
		config.IrodsSslSaltSize,
		config.IrodsHashRounds,
	)
	if err != nil {
		return nil, err
	}

	account.SSLConfiguration = sslConfig
	account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
	account.ClientServerNegotiation = true

	return account, nil
}

func NewFileSystem(appName string, config Config) (*fs.FileSystem, error) {
	acct, err := NewAccount(config)
	if err != nil {
		return nil, err
	}

	fsys, err := fs.NewFileSystemWithDefault(acct, appName)
	if err != nil {
		return nil, err
	}

	return fsys, nil
}
