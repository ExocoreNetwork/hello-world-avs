package avs

import (
	"context"
	chain "github.com/ExocoreNetwork/exocore-avs/core/chainio"
	"github.com/ExocoreNetwork/exocore-avs/core/chainio/eth"
	"github.com/ExocoreNetwork/exocore-avs/types"
	"github.com/ExocoreNetwork/exocore-sdk/chainio/txmgr"
	"github.com/ExocoreNetwork/exocore-sdk/logging"
	sdklogging "github.com/ExocoreNetwork/exocore-sdk/logging"
	"github.com/ExocoreNetwork/exocore-sdk/signerv2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/rand"
	"os"
	"time"
)

const (
	avsName = "hello-avs-demo"
)

type Avs struct {
	logger    logging.Logger
	avsWriter chain.ExoWriter
	avsReader chain.ExoReader
}

// NewAvs creates a new Avs with the provided config.

func NewAvs(c *types.NodeConfig) (*Avs, error) {
	var logLevel sdklogging.LogLevel
	if c.Production {
		logLevel = sdklogging.Production
	} else {
		logLevel = sdklogging.Development
	}
	logger, err := sdklogging.NewZapLogger(logLevel)
	if err != nil {
		return nil, err
	}

	ethRpcClient, err := eth.NewClient(c.EthRpcUrl)
	if err != nil {
		logger.Error("Cannot create http ethclient", "err", err)
		return nil, err
	}
	chainId, err := ethRpcClient.ChainID(context.Background())
	if err != nil {
		logger.Error("Cannot get chainId", "err", err)
		return nil, err
	}

	ecdsaKeyPassword, ok := os.LookupEnv("OPERATOR_ECDSA_KEY_PASSWORD")
	if !ok {
		logger.Info("OPERATOR_ECDSA_KEY_PASSWORD env var not set. using empty string")
	}

	signerV2, _, err := signerv2.SignerFromConfig(signerv2.Config{
		KeystorePath: c.AVSEcdsaPrivateKeyStorePath,
		Password:     ecdsaKeyPassword,
	}, chainId)
	if err != nil {
		panic(err)
	}

	txMgr := txmgr.NewSimpleTxManager(ethRpcClient, logger, signerV2, common.HexToAddress(c.AVSOwnerAddress))
	avsWriter, err := chain.BuildExoChainWriter(
		common.HexToAddress(c.AVSAddress),
		ethRpcClient,
		logger,
		txMgr)
	if err != nil {
		logger.Error("Cannot create avsWriter", "err", err)
		return nil, err
	}

	avsReader, err := chain.BuildExoChainReader(
		common.HexToAddress(c.AVSAddress),
		ethRpcClient,
		logger)
	if err != nil {
		logger.Error("Cannot create exoChainReader", "err", err)
		return nil, err
	}
	info, err := avsReader.GetAVSInfo(&bind.CallOpts{}, c.AVSAddress)
	if err != nil {
		logger.Error("Cannot GetAVSInfo", "err", err)
		return nil, err
	}
	if info == "" {
		_, err = avsWriter.RegisterAVSToExocore(context.Background(),
			avsName,
			c.MinStakeAmount,
			common.HexToAddress(c.AVSAddress),
			common.HexToAddress(c.AVSRewardAddress),
			common.HexToAddress(c.AVSSlashAddress),
			c.AvsOwnerAddresses,
			c.AssetIds,
			c.AvsUnbondingPeriod,
			c.MinSelfDelegation,
			c.EpochIdentifier,
			c.Params,
		)
		if err != nil {
			logger.Error("register Avs failed ", "err", err)
			return &Avs{}, err
		}
	}

	return &Avs{
		logger:    logger,
		avsWriter: avsWriter,
		avsReader: avsReader,
	}, nil
}

func (avs *Avs) Start(ctx context.Context) error {
	avs.logger.Infof("Starting avs.")
	ticker := time.NewTicker(50 * time.Second)
	avs.logger.Infof("Avs owner set to send new task every 50 seconds...")
	defer ticker.Stop()
	taskNum := int64(1)
	// send the first task
	_ = avs.sendNewTask()
	taskNum++
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			avs.logger.Info("sendNewTask-num:", taskNum)
			err := avs.sendNewTask()
			taskNum++
			if err != nil {
				// we log the errors inside sendNewTask() so here we just continue to the next task
				continue
			}
		}
	}
}

// sendNewTask sends a new task to the task manager contract.
func (avs *Avs) sendNewTask() error {
	avs.logger.Info("Avs sending new task")
	_, err := avs.avsWriter.CreateNewTask(
		context.Background(),
		GenerateRandomName(5),
		types.TaskResponsePeriod,
		types.TaskChallengePeriod,
		types.ThresholdPercentage,
		types.TaskStatisticalPeriod)

	if err != nil {
		avs.logger.Error("Avs failed to sendNewTask", "err", err)
		return err
	}

	return nil
}
func GenerateRandomName(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
