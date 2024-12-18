package config

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/ExocoreNetwork/exocore-avs/core/chainio/eth"
	"github.com/ExocoreNetwork/exocore-sdk/chainio/txmgr"
	sdklogging "github.com/ExocoreNetwork/exocore-sdk/logging"
	"github.com/ExocoreNetwork/exocore-sdk/signerv2"

	blscommon "github.com/prysmaticlabs/prysm/v4/crypto/bls/common"
)

type Config struct {
	Production      bool `yaml:"production"`
	EcdsaPrivateKey *ecdsa.PrivateKey
	BlsPrivateKey   *blscommon.SecretKey
	Logger          sdklogging.Logger
	// we need the url for the exocore-sdk currently... eventually standardize api to
	// only take an ethclient or an rpcUrl (and build the ethclient at each constructor site)
	EthHttpRpcUrl              string
	EthWsRpcUrl                string
	EthHttpClient              eth.EthClient
	EthWsClient                eth.EthClient
	OperatorStateRetrieverAddr common.Address
	AvsRegistryCoordinatorAddr common.Address
	AggregatorServerIpPortAddr string
	RegisterOperatorOnStartup  bool
	// json:"-" skips this field when marshaling (only used for logging to stdout), since SignerFn doesnt implement marshalJson
	SignerFn                 signerv2.SignerFn `json:"-"`
	TxMgr                    txmgr.TxManager
	AggregatorAddress        common.Address
	EcdsaPrivateKeyStorePath string `yaml:"ecdsa_private_key_store_path"`
}

var (
	FileFlag = cli.StringFlag{
		Name:     "config",
		Required: true,
		Usage:    "Load configuration from `FILE`",
	}

	EcdsaPrivateKeyFlag = cli.StringFlag{
		Name:     "ecdsa-private-key",
		Usage:    "Ethereum private key",
		Required: true,
		EnvVar:   "ECDSA_PRIVATE_KEY",
	}
	/* Optional Flags */
)

var requiredFlags = []cli.Flag{
	FileFlag,
	EcdsaPrivateKeyFlag,
}

var optionalFlags = []cli.Flag{}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag
